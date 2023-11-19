

function parseArguments(){
    const args = process.argv
    let capturedArgMap = {}
    let targetDir = ""
    let captureNextAs = "target"
    for(let idx in args ){
        const argument = args[idx]
        if( captureNextAs ){
            capturedArgMap[captureNextAs] = argument
            captureNextAs = ''
        } else  if( argument === "-t" ){
            targetDir = args[idx+1]
            captureNextAs = 'target'
        } else {

        }
    }

    return { cmd: '', paramMap:capturedArgMap}
}
const parsedArgs = parseArguments()

console.log(parsedArgs.paramMap.target)

const fs = require('fs')
const path = require("path");

const packageJsonPath = path.join(parsedArgs.paramMap.target, "package.json")
if( fs.existsSync(packageJsonPath) ){
    const packageJsonData = fs.readFileSync(packageJsonPath)
    const packageJson = JSON.parse(packageJsonData)

    fs.writeFileSync(packageJsonPath +'.ori', packageJsonData)

    // publishConfig field 제거
    delete packageJson.publishConfig

    fs.writeFileSync(packageJsonPath, JSON.stringify(packageJson))
}