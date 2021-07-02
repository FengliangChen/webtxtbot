/*
 @Author      : Simon Chen
 @Email       : bafelem@gmail.com
 @datetime    : 2021-07-02 15:23:58
 @Description : Description
 @FileName    : t3.go
*/

package finalartwork

import (
    "io/ioutil"
    "log"
    "os"
    "os/exec"
)

var baseScript = `
osascript <<EOF
set theRecipients to {{name:"Faris", email:"faris.abukwiek@walmart.com"}}
set ccRecpients to {{name:"service", email:"service@benchmarkdesign.org"}, {name:"marika", email:"marika@benchmarkdesign.org"},{name:"chris", email:"chris@benchmarkdesign.org"}}
tell application "Mail"
    set theMessage to make new outgoing message with properties {subject:"$etitle", content:"$econtent", visible:true}
    
    repeat with theRecipient in theRecipients
        tell theMessage
            make new to recipient at end of to recipients with properties {name:name of theRecipient, address:email of theRecipient}
        end tell
    end repeat
    repeat with ccRecipient in ccRecpients
        tell theMessage
            make new cc recipient at end of cc recipients with properties {name:name of ccRecipient, address:email of ccRecipient}
        end tell
    end repeat
end tell
EOF
`

func MakeEmail(title string, content string) {

    // Create our Temp File:  This will create a filename like /tmp/prefix-123456
    // We can use a pattern of "pre-*.txt" to get an extension like: /tmp/pre-123456.txt
    tmpFile, err := ioutil.TempFile(os.TempDir(), "prefix-*.sh")
    if err != nil {
        log.Fatal("Cannot create temporary file", err)
    }

    // Remember to clean up the file afterwards
    defer os.Remove(tmpFile.Name())

    script := genScript(title, content)
    byte_script := []byte(script)
    if _, err = tmpFile.Write(byte_script); err != nil {
        log.Fatal("Failed to write to temporary file", err)
    }

    // Close the file
    if err := tmpFile.Close(); err != nil {
        log.Fatal(err)
    }

    err = RunScript(tmpFile.Name())
    if err != nil {
        log.Fatal(err)
    }

}

func RunScript(path string) error {
    cmd := exec.Command("bash", path)
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}

func genScript(title string, content string) string {
    firstLine := "etitle=" + "\"" + title + "\"" + "\n"
    secondLine := "econtent=" + "\"" + content + "\"" + "\n"

    script := firstLine + secondLine + baseScript
    return script
}
