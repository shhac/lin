import { Command } from "commander";
import { getPackageVersion } from "./lib/version.ts";
import { registerAuthCommand } from "./cli/auth/index.ts";
import { registerCycleCommand } from "./cli/cycle/index.ts";
import { registerIssueCommand } from "./cli/issue/index.ts";
import { registerLabelCommand } from "./cli/label/index.ts";
import { registerProjectCommand } from "./cli/project/index.ts";
import { registerTeamCommand } from "./cli/team/index.ts";
import { registerUsageCommand } from "./cli/usage-command.ts";
import { registerUserCommand } from "./cli/user/index.ts";

const program = new Command();
program.name("lin").description("Linear CLI for humans and LLMs").version(getPackageVersion());

registerAuthCommand({ program });
registerProjectCommand({ program });
registerIssueCommand({ program });
registerTeamCommand({ program });
registerUserCommand({ program });
registerLabelCommand({ program });
registerCycleCommand({ program });
registerUsageCommand({ program });

program.parse(process.argv);
if (!process.argv.slice(2).length) {
  program.outputHelp();
}
