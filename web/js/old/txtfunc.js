var obj
var rvst = true

function rvstCheck(){
	if (document.getElementById("rvstCheck").checked == true){
		return false
	}else{
		return true
	}
}

function validateForm(){
	obj = null

	var x = document.getElementById("job").value
	if (x == null || x == ""){
		alert("需要输入单号。");
		return false;
	}
	rvst = rvstCheck()
	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			try {
				obj = JSON.parse(txt);
			}
			catch(err) {
				DelePara()
				var element=document.getElementById("txtdv");
				var para = document.createElement("p")
				para.innerText = xmlhttp.responseText;
				element.appendChild(para)
				return
			}
			if (rvst) {
				rvstOutput()
			}else{
				NormalTxtOutput()
			}
		}
	}
	if (rvst) {
		xmlhttp.open("GET","/draftArtworkText?j="+x+"&rvst=true",true);
	}else{
		xmlhttp.open("GET","/draftArtworkText?j="+x,true);
	}
	
	xmlhttp.send();

}

function ProofTxtOutput(){
	if (obj == null) {
		return
	}
	DelePara()
	var element=document.getElementById("txtdv");

	var lineBreak = document.createElement("br")
	var welcomeInput = document.createElement("input")
	welcomeInput.setAttribute("id", "welcomeInput")
	welcomeInput.setAttribute("value", "")
	welcomeInput.setAttribute("placeholder", "Change Supplier")
	welcomeInput.setAttribute("oninput", "hiFunction()")
	element.appendChild(welcomeInput)
	element.appendChild(lineBreak)

	for (var i = 0; i < obj.TxtBodies.length; i++ ){
		var BodyStruct = obj.TxtBodies[i];
		var para = document.createElement("p")
		var btn = document.createElement("button")
		var txtspan = document.createElement("span")
		var lineBreak = document.createElement("br")
		var cpcommand = "copyElementToClipboard(" + "document.getElementById(\"p" + i + "\"));" 
		txtspan.innerText = BodyStruct.TxtBody
		para.innerHTML= ProofTitle(BodyStruct.TxtCount) + "<span>" + txtspan.innerHTML + "</span>"
		para.setAttribute("id", "p" + i)
		para.setAttribute("contenteditable", "true")
		btn.innerText = "Copy"
		btn.setAttribute('onclick', cpcommand);
		element.appendChild(para)
		element.appendChild(btn)
		element.appendChild(lineBreak)
		element.appendChild(lineBreak.cloneNode())
	}

}

function rvstOutput(){
	if (obj == null) {
		return
	}
	DelePara()

	var element=document.getElementById("txtdv");

	var lineBreak = document.createElement("br")

	var rvst1Input = document.createElement("input")
	rvst1Input.setAttribute("id", "rvst1Input")
	rvst1Input.setAttribute("value", "")
	rvst1Input.setAttribute("placeholder", "egg.translation")
	rvst1Input.setAttribute("oninput", "rvstFunction('rvst1Input','updateReason1')")
	element.appendChild(rvst1Input)

	var rvst2Input = document.createElement("input")
	rvst2Input.setAttribute("id", "rvst2Input")
	rvst2Input.setAttribute("value", "")
	rvst2Input.setAttribute("placeholder", "egg. final approval")
	rvst2Input.setAttribute("oninput", "rvstFunction('rvst2Input','updateReason2')")
	element.appendChild(rvst2Input)

	element.appendChild(lineBreak)

	for (var i = 0; i < obj.TxtBodies.length; i++ ){
		var BodyStruct = obj.TxtBodies[i];
		var para = document.createElement("p")
		var btn = document.createElement("button")
		var txtspan = document.createElement("span")
		var lineBreak = document.createElement("br")
		var cpcommand = "copyElementToClipboard(" + "document.getElementById(\"p" + i + "\"));" 
		txtspan.innerText = BodyStruct.TxtBody
		para.innerHTML= rvstTitle(BodyStruct.TxtCount) + "<span>" + txtspan.innerHTML + "</span>"
		para.setAttribute("id", "p" + i)
		para.setAttribute("contenteditable", "true")
		btn.innerText = "Copy"
		btn.setAttribute('onclick', cpcommand);
		element.appendChild(para)
		element.appendChild(btn)
		element.appendChild(lineBreak)
		element.appendChild(lineBreak.cloneNode())
	}

}

function rvstTitle(count){
	 if (count == 1) {
    var str=`
<span>
<span>The following </span>
<span name="updateReason1" style="color: rgb(255, 0, 0);">revised file </span>
<span>is for your </span>
<span name="updateReason2" style="color: rgb(255, 0, 0);">approval:</span>
<br>
<br>
</span>
`
	return str
	}else{
    var str = `
<span>
<span>The following </span>
<span style="color: rgb(255, 0, 0);">Count </span>
<span name="updateReason1" style="color: rgb(255, 0, 0);">revised files </span>
<span>are for your </span>
<span name="updateReason2" style="color: rgb(255, 0, 0);">approval:</span>
<br>
<br>
</span>
`
str = str.replace("Count", count)
return str
	}

}

function ProofTitle(count) {
  if (count == 1) {
    var Str=`
<span>
<span name="welcome" onclick="changeProofer(this)">Hi Ella,</span>
<br>
<br>
<span>The following </span>
<span style="color: rgb(255, 0, 0);">SUPPLIER CREATED FILE </span>
<span>is for your </span>
<span style="color: rgb(255, 0, 0);">approval:</span>
<br>
<br>
</span>
`
return Str
  }else{
    var Str = `
<span>
<span name="welcome" onclick="changeProofer(this)">Hi Ella,</span>
<br>
<br>
<span>The following </span>
<span style="color: rgb(255, 0, 0);">Count </span>
<span style="color: rgb(255, 0, 0);">SUPPLIER CREATED FILES </span>
<span>are for your </span>
<span style="color: rgb(255, 0, 0);">approval:</span>
<br>
<br>
</span>
`
Str = Str.replace("Count", count)
return Str
  }

}

function changeProofer(obj){
	var name = obj.innerText

	switch (name){
		case "Hi Ella,":
		obj.innerText = "Hi Sam,"
		break;
	}
}

function NormalTxtOutput() {
		if (obj == null) {
		return
	}
	DelePara()
	var element=document.getElementById("txtdv");
	var phqPara = document.createElement("p")
	var phqbtn = document.createElement("button")
	var lineBreak = document.createElement("br")
	var cpcommand = "copy(" + "document.getElementById(\"phqtxt\").innerText);" 

	phqPara.innerText = obj.PHQ
	phqPara.setAttribute("id", "phqtxt")
	phqPara.setAttribute("contenteditable", "true")
	element.appendChild(phqPara)
	phqbtn.innerText = "Copy"
	phqbtn.setAttribute('onclick', cpcommand);
	element.appendChild(phqbtn)
	element.appendChild(lineBreak)
	element.appendChild(lineBreak.cloneNode())
	element.appendChild(lineBreak.cloneNode())

	var welcomeInput = document.createElement("input")
	welcomeInput.setAttribute("id", "welcomeInput")
	welcomeInput.setAttribute("value", "")
	welcomeInput.setAttribute("placeholder", "Change Supplier")
	welcomeInput.setAttribute("oninput", "hiFunction()")
	element.appendChild(welcomeInput)
	element.appendChild(lineBreak)

	for (var i = 0; i < obj.TxtBodies.length; i++ ){
		var BodyStruct = obj.TxtBodies[i];
		var para = document.createElement("p")
		var btn = document.createElement("button")
		var txtspan = document.createElement("span")
		var lineBreak = document.createElement("br")
		var cpcommand = "copyElementToClipboard(" + "document.getElementById(\"p" + i + "\"));"
		txtspan.innerText = BodyStruct.TxtBody
		para.innerHTML= NormalTitle(BodyStruct.TxtCount) + "<span>" + txtspan.innerHTML + "</span>"
		para.setAttribute("id", "p" + i)
		para.setAttribute("contenteditable", "true")
		btn.innerText = "Copy"
		btn.setAttribute('onclick', cpcommand);
		element.appendChild(para)
		element.appendChild(btn)
		element.appendChild(lineBreak)
		element.appendChild(lineBreak.cloneNode())
	}
}

function NormalTitle(count) {
	var str
	if (count == 1) {
		str = `
<span>
<span name="welcome">Hi Supplier,</span>
<br>
<br>
<span>The following </span>
<span style="color: rgb(255, 0, 0);">draft artwork file </span>
<span>is for your </span>
<span style="color: rgb(255, 0, 0);">first approval:</span>
<br>
<br>
</span>
`
		return str
	}else{
		str = `
<span>
<span name="welcome">Hi Supplier,</span>
<br>
<br>
<span>The following </span>
<span style="color: rgb(255, 0, 0);">Count draft artwork files </span>
<span>are for your </span>
<span style="color: rgb(255, 0, 0);">first approval:</span>
<br>
<br>
</span>
`
		str = str.replace("Count", count)
		return str
	}
}

function copy(text) {
    var input = document.createElement('textarea');
    input.innerHTML = text;
    document.body.appendChild(input);
    input.select();
    var result = document.execCommand('copy');
    document.body.removeChild(input);
    return result;
}

function DelePara(){
	var div = document.getElementById("txtdv")
	div.innerHTML = ""
}

function clickbackgroundColor(btn){
	btn.style.backgroundColor = "#888888"
}

function copyElementToClipboard(element) {
  window.getSelection().removeAllRanges();
  let range = document.createRange();
  range.selectNode(typeof element === 'string' ? document.getElementById(element) : element);
  window.getSelection().addRange(range);
  document.execCommand('copy');
  window.getSelection().removeAllRanges();
}

function hiFunction(){
	var x = document.getElementById('welcomeInput').value
	var elements = document.getElementsByName("welcome")
	x = "Hi " + x + ","

	for(var i=0; i<elements.length; i++){
		elements[i].innerText = x;
	}
}

function rvstFunction(iputid,nametag){
	var x = document.getElementById(iputid).value
	var elements = document.getElementsByName(nametag)

	switch (nametag){
		case "updateReason1":
		x = x + " "
		break;
		case "updateReason2":
		x = x + "approval:"
		break;

	}
	for(var i=0; i<elements.length; i++){
		elements[i].innerText = x;
	}
}

function FinalProofTxtOutput(){
	if (obj == null) {
		return
	}
	DelePara()
	var element=document.getElementById("txtdv");

	var lineBreak = document.createElement("br")
	var welcomeInput = document.createElement("input")
	welcomeInput.setAttribute("id", "welcomeInput")
	welcomeInput.setAttribute("value", "")
	welcomeInput.setAttribute("placeholder", "Change Supplier")
	welcomeInput.setAttribute("oninput", "hiFunction()")
	element.appendChild(welcomeInput)
	element.appendChild(lineBreak)

	for (var i = 0; i < obj.TxtBodies.length; i++ ){
		var BodyStruct = obj.TxtBodies[i];
		var para = document.createElement("p")
		var btn = document.createElement("button")
		var txtspan = document.createElement("span")
		var lineBreak = document.createElement("br")
		var cpcommand = "copyElementToClipboard(" + "document.getElementById(\"p" + i + "\"));" 
		txtspan.innerText = BodyStruct.TxtBody
		para.innerHTML= FinalProofTitle(BodyStruct.TxtCount) + "<span>" + txtspan.innerHTML + "</span>"
		para.appendChild(FinalProofTail())
		para.setAttribute("id", "p" + i)
		para.setAttribute("contenteditable", "true")
		btn.innerText = "Copy"
		btn.setAttribute('onclick', cpcommand);
		element.appendChild(para)
		element.appendChild(btn)
		element.appendChild(lineBreak)
		element.appendChild(lineBreak.cloneNode())
	}

}

function FinalProofTitle(count){
	if (count == 1){
		var str=
`<span>
<span name="welcome">Hi Supplier,</span>
<br>
<br>
<span>The following file is approved by Walmart and our side, please proceed.</span>
<br>
<br>
</span>`
		return str
	}else{
		var str = 
`<span>
<span name="welcome">Hi Supplier,</span>
<br>
<br>
<span>The following files are approved by Walmart and our side, please proceed.</span>
<br>
<br>
</span>`
		return str
	}

}


function FinalProofTail(){
	var str=`Please note: Our approval of the files is for the graphic consistency with Walmart’s style guide/modular look! You still need to send artwork to the testing agency for approval for any legal/compliance issues. Our approval does not mean the file is approved as far as compliance for Warnings, Dimensions, Ingredients, or any other legal requirement. Our approval only means it meets the graphic requirements provided by Walmart.
	`
	var tail = document.createElement("span")
	tail.innerText = str
	tail.setAttribute("style","color: rgb(255, 0, 0)")
	return tail
}

