import { Command } from "commander";
import { getPackageVersion } from "./lib/version.ts";
import { configureTruncation } from "./lib/truncation.ts";
import { getSettings } from "./lib/config.ts";
import { registerAuthCommand } from "./cli/auth/index.ts";
import { registerCycleCommand } from "./cli/cycle/index.ts";
import { registerDocumentCommand } from "./cli/document/index.ts";
import { registerIssueCommand } from "./cli/issue/index.ts";
import { registerLabelCommand } from "./cli/label/index.ts";
import { registerProjectCommand } from "./cli/project/index.ts";
import { registerRoadmapCommand } from "./cli/roadmap/index.ts";
import { registerTeamCommand } from "./cli/team/index.ts";
import { registerConfigCommand } from "./cli/config-command.ts";
import { registerUsageCommand } from "./cli/usage-command.ts";
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

registerAuthCommand({ program });
registerProjectCommand({ program });
registerRoadmapCommand({ program });
registerDocumentCommand({ program });
registerIssueCommand({ program });
registerTeamCommand({ program });
registerUserCommand({ program });
registerLabelCommand({ program });
registerCycleCommand({ program });
registerConfigCommand({ program });
registerUsageCommand({ program });

program.parse(process.argv);
if (!process.argv.slice(2).length) {
  program.outputHelp();
}
