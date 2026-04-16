import { Command } from "commander";
import { getPackageVersion } from "./lib/version.ts";
import { configureTruncation } from "./lib/truncation.ts";
import { getSettings } from "./lib/config.ts";
import { registerApiCommand } from "./cli/api/index.ts";
import { registerAuthCommand } from "./cli/auth/index.ts";
import { registerCycleCommand } from "./cli/cycle/index.ts";
import { registerDocumentCommand } from "./cli/document/index.ts";
import { registerFileCommand } from "./cli/file/index.ts";
import { registerIssueCommand } from "./cli/issue/index.ts";
import { registerLabelCommand } from "./cli/label/index.ts";
import { registerProjectCommand } from "./cli/project/index.ts";
import { registerRoadmapCommand } from "./cli/roadmap/index.ts";
import { registerTeamCommand } from "./cli/team/index.ts";
import { registerConfigCommand } from "./cli/config/index.ts";
import { registerUsageCommand } from "./cli/usage/index.ts";
import { registerUserCommand } from "./cli/user/index.ts";

const program = new Command();
program.name("lin").description("Linear CLI for humans and LLMs").version(getPackageVersion());
program.option(
  "--expand <fields>",
  "Expand truncated fields (comma-separated: description,body,content)",
);
program.option("--full", "Show full content for all truncated fields");
program.hook("preAction", (thisCommand) => {
  const opts = thisCommand.opts();
  const settings = getSettings();
  configureTruncation({
    expand: opts.expand,
    full: opts.full,
    maxLength: settings.truncation?.maxLength,
  });
});

registerApiCommand({ program });
registerAuthCommand({ program });
registerProjectCommand({ program });
registerRoadmapCommand({ program });
registerDocumentCommand({ program });
registerFileCommand({ program });
registerIssueCommand({ program });
registerTeamCommand({ program });
registerUserCommand({ program });
registerLabelCommand({ program });
registerCycleCommand({ program });
registerConfigCommand({ program });
registerUsageCommand({ program });

const HELP_HINT =
  "\nRun 'lin usage' for detailed docs. Run 'lin <command> usage' for per-command details.";
const SUBCMD_HELP_HINT = (name: string) =>
  `\nRun 'lin ${name} usage' for detailed docs including all fields, valid values, and examples.`;

program.addHelpText("after", HELP_HINT);
for (const cmd of program.commands) {
  if (cmd.commands.some((c) => c.name() === "usage")) {
    cmd.addHelpText("after", SUBCMD_HELP_HINT(cmd.name()));
  }
}

program.parse(process.argv);
if (!process.argv.slice(2).length) {
  program.outputHelp();
}
