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