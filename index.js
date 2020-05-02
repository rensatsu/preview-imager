const fs = require("fs");
const path = require("path");
const spawn = require("child_process").spawn;

class PreviewImagerError extends Error {
    constructor(message, output = null) {
        super(message);
        this.name = "PreviewImagerError";
        this.output = output;
    }
}

let cmd = null;
let cmdCwd = null;
const cmdParams = [];

if (fs.existsSync(path.join(__dirname, "/dist/preview-imager"))) {
    cmd = path.join(__dirname, "/dist/preview-imager");
    cmdCwd = __dirname;
} else if (fs.existsSync(path.join(__dirname, "/go/app.go"))) {
    cmd = "go";
    cmdParams.push("run", path.join(__dirname, "/go/app.go"));
    cmdCwd = path.join(__dirname, "/go/");
}

class PreviewImager {
    constructor(output = null, title = null) {
        if (cmd === null) {
            throw new PreviewImagerError("Unable to find preview-imager executable");
        }

        if (output === null) throw new PreviewImagerError("Output file is not defined");
        if (title === null) throw new PreviewImagerError("Title is not defined");

        this._options = {
            targetPath: output,
            title: title,
            colorBg: null,
            colorBgOverlay: null,
            colorFg: null,
            siteName: null,
            imagePath: null,
            logoPath: null,
            logoX: null,
            logoY: null,
            width: null,
            height: null,
            paddingX: null,
            paddingY: null,
            lineSpacing: null,
            fontSize: null,
            fontSizeSite: null,
            paddingYSite: null,
            blurStrength: null,
            fontTitle: null,
            fontSiteName: null,
        };
    }

    set(option, value) {
        if (!(option in this._options)) return;
        this._options[option] = value;
    }

    generate() {
        return new Promise((resolve, reject) => {
            if (this.title === null) {
                reject(new PreviewImagerError("Title is not defined"));
                return;
            }

            const optArray = Object.entries(this._options).reduce((acc, [k, v]) => {
                if (v !== null) {
                    acc.push(`--${k}`, v);
                }

                return acc;
            }, []);

            const cmdOpt = [...cmdParams, ...optArray];
            const goChild = spawn(cmd, cmdOpt);

            let shOutput = [];

            goChild.stdout.on("data", (data) => shOutput.push(`stdout: ${data}`));
            goChild.stderr.on("data", (data) => shOutput.push(`stderr: ${data}`));

            goChild.on("exit", (code) => {
                if (code === 0) {
                    resolve(shOutput.join("\n"));
                } else {
                    reject(new PreviewImagerError(
                        "Execution failed",
                        shOutput.join("\n")
                    ));
                    return;
                }
            });
        });
    }
}

module.exports = { PreviewImager, PreviewImagerError };
