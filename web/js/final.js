var finalartworkObj;
var previousInput;
var finalRecordObj;
var finishCompress = false;

function finalProcessByKey(event) {
	var x = event.key;
	if (x == "Enter"){
		FinalArtwork();
	}
}

function FinalArtwork(){
	var job = document.getElementById("finalArtworkInput").value;
	if (job == null || job == ""){
	alert("需要输入单号。");
	return false;
	}
	if (job.length != 6){
	alert("需要输入6位数的单号。");
	return false;
	}
	if (previousInput == job){
		alert("与上一次输入重复！");
		return;
	}
	previousInput = job;
	ClearContent("mainFinal")
	ClearContent("compressOption")
	ClearContent("processStateBar")

	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			try {
				finalartworkObj = JSON.parse(txt);
			}
			catch(err) {
				var element=document.getElementById("mainFinal");
				var para = document.createElement("p")
				para.innerText = xmlhttp.responseText;
				element.appendChild(para)
				return
			}
			FinalBodyOutput()
			CompressOptionBuild()
			autoEmailButton()
		}
	}

	xmlhttp.open("GET","/query?j="+job,true);
	xmlhttp.send();
}

function FinalBodyOutput(){
	if (finalartworkObj == null){
		return
	}
	ClearContent("mainFinal")

	var outPutTxt = finalartworkObj.Body
	var mainDiv = document.getElementById("mainFinal");
	var para = document.createElement("p");
	if (outPutTxt.length === 0){
		para.innerHTML = outPutTxt = finalartworkObj.Error
		mainDiv.appendChild(para)
		return
	}
	var inputItem = document.createElement("input")
	inputItem.setAttribute("placeholder", "更新 Supplier")
	inputItem.setAttribute("oninput", "onchangeSupplier(this.value)")
	mainDiv.appendChild(inputItem)


	para.innerHTML = BodyInnerHTMLpart()
	para.setAttribute("id", "mainFinalText")
	mainDiv.appendChild(para)

	var btn = document.createElement("button");
	btn.innerText = "Copy"
	btn.setAttribute('onclick', "CopyFinalText()");
	mainDiv.appendChild(btn)

}

function CopyFinalText(){
	var mainDiv = document.getElementById("mainFinalText");
	copyElementToClipboard(mainDiv)
}

function copyElementToClipboard(element) {
  window.getSelection().removeAllRanges();
  let range = document.createRange();
  range.selectNode(typeof element === 'string' ? document.getElementById(element) : element);
  window.getSelection().addRange(range);
  document.execCommand('copy');
  window.getSelection().removeAllRanges();
}


function ClearContent(itemID){
	var divObj = document.getElementById(itemID);
	divObj.innerHTML = "";
}

function CompressOptionBuild(){
	if (finalartworkObj.Zippable == true ){
		var mainDiv = document.getElementById("compressOption");
		var para = document.createElement("p");
		var span = document.createElement("span");
		span.innerText = "工单号：" + finalartworkObj.Jobcode + " 文件大小："+ (finalartworkObj.TotalFileSize/1024/1024).toFixed(2) + " MB";
		para.appendChild(span);
		mainDiv.appendChild(para)
		var btn = document.createElement("button");
		btn.innerText = "压缩"
		btn.setAttribute("name", finalartworkObj.Token);
		btn.setAttribute('onclick', "Compress(this.name);this.disabled='true'");
		mainDiv.appendChild(btn);
	}
}

function Compress(itemId){
	var xmlhttp;
	xmlhttp=new XMLHttpRequest();

	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var para = document.createElement("p");
			var mainDiv = document.getElementById("compressOption");
			var txt = xmlhttp.responseText
			para.innerText = txt
			para.setAttribute("id", "zip"+itemId)
			para.style.fontWeight = "bold"
			mainDiv.appendChild(para)
		}
	}
	xmlhttp.open("GET","/compress?j="+itemId, true);
	xmlhttp.send();
}

var trackObj
var tracking = setInterval(TrackRequest, 1000);

function TrackRequest(){
	var xmlhttp;
	xmlhttp=new XMLHttpRequest();

	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			if (txt == "null"){
				return
			}
			
			try {
				trackObj = JSON.parse(txt);
			}
			catch(err) {
				return
			}

			if (trackObj != null){
				TrackBarUpdate()
			}
		}
	}
	xmlhttp.open("GET","/track", true);
	xmlhttp.send();
}

function TrackBarUpdate(){
	for(var i = 0; i < trackObj.length; i++ )
	{

		var objtxt = document.getElementById("sub" + trackObj[i].Token);
		var barobj = document.getElementById(trackObj[i].Token);
		var percentage = (100*trackObj[i].CompressedSize / trackObj[i].TotalFileSize).toFixed(0);

		DynamicTitle(trackObj[i].Jobcode, percentage);

		if (barobj == null) {
			BuildBars(trackObj[i], percentage);
			if (percentage == 100){
				UpdateZipBtnTxt(trackObj[i].Token)
			}
			continue
		}else{
			objtxt.innerText = ItemStatusInfo(trackObj[i]);
			UpdateBar(trackObj[i].Token, percentage);

			if (percentage == 100){
				UpdateZipBtnTxt(trackObj[i].Token)
			}
		}

	}
}

function UpdateZipBtnTxt(itemID){
	var element = document.getElementById("zip"+itemID)
	if (element == null){
		return
	}
	if (finishCompress == true){
		return
	}
	element.innerText = "已完成！"
	alert("压缩文件完成！")
	finishCompress = true;
}


function ItemStatusInfo(itemObj){
	var name = itemObj.Jobcode
	var size = itemObj.CompressedSize
	var totalsize = itemObj.TotalFileSize
	size = (size/1024/1024).toFixed(2)
	totalsize = (totalsize/1024/1024).toFixed(2)
	var txt = "压缩文件名：" + name +" 总大小：" + totalsize+"MB 已压缩："+size + "MB"
	return txt
}

function DynamicTitle(jobcode, percentage){
	document.title = percentage + "% - " + "Zipping..."
	if (percentage == 100) {
		document.title = "Final artwork"
	}
}

function onchangeSupplier(onchangeText){
	var item = document.getElementById("supplierLine");
	item.innerText = "Hi " + onchangeText + ",";

}

function BuildBars(itemObj,percentage){
	var barID = itemObj.Token;
	var subID = "sub" + barID


	var barDiv = document.createElement("div")
	var subBarDiv = document.createElement("div")

	barDiv.setAttribute("class", "w3-light-grey")
	subBarDiv.setAttribute("class", "w3-container w3-green")
	subBarDiv.setAttribute("style", "height:24px;width:" + percentage + "%");
	subBarDiv.innerHTML = percentage + "%";
	subBarDiv.setAttribute("id", barID)
	barDiv.appendChild(subBarDiv)



	var mainDiv = document.getElementById("processStateBar");
	var para = document.createElement("p");
	para.innerText = ItemStatusInfo(itemObj)
	para.setAttribute("id", subID)


	mainDiv.insertBefore(barDiv, mainDiv.firstChild);
	mainDiv.insertBefore(para, mainDiv.firstChild);
}

function UpdateBar(itemID, pencentage){
	var item = document.getElementById(itemID);
	item.style.width = pencentage + '%';
	item.innerHTML = pencentage + '%';

}


function BodyInnerHTMLpart(){

	var body = document.createElement("span")
	body.innerText = finalartworkObj.Body



	var head = `
<span id="supplierLine">Hi Supplier,</span><br><br>

<span>Here are the final files for following items, please confirm receipt.</span>
<br><br>
<span id="linkClaim1">Link:<br><br></span>
`
	var head2 = `
<span id="supplierLine">Hi Supplier,</span><br><br>

<span>Here is final file for following item, please confirm receipt.</span>
<br><br>
<span id="linkClaim1">Link:<br><br></span>
`




	var tail = `
<br>
<span>We have provided the files in 2 formats for you (Illustrator files, and PDF files). The artwork is exactly the same in each format, BUT they are for different purposes.</span><br><br>
<span>Illustrator format (AI files) --- these are for PRINTING ONLY - please send these to your PRINTER. </span>
<span style="color: rgb(241, 6, 233);">Please use Adobe Illustrator CC 2018 or higher version when open it, if you need lower version, please feel free to contact us. </span><br><br>
 
<span>PDF format --- these are for your visual REFERENCE ONLY. These files are locked, so that no changes can be made. Do not mass print from the PDFs!</span><br>
<span>If you are looking for WMT approved printers, please feel free to contact us for a printing quote.</span><br><br>
<span>If you have any questions, or problems with the files, please don't hesitate to let us know, we’ll be happy to help.</span><br><br>

<span id="linkClaim" style="color: rgb(255, 0, 0)" onclick="DeleteLinkText()">IMPORTANT: The above download link will be available for 7 days. We will delete the files from the server at that time. If you do not download the files within 7 days, please contact us and we will re-upload to the server. There will be a 1 hour service fee to do this.</span>
	`
if (filesPlural(finalartworkObj.Body)){
	return head + body.innerHTML + tail
}else{
	return head2 + body.innerHTML + tail
}
}

function DeleteLinkText(){
	var link = document.getElementById("linkClaim1")
	var linktxt = document.getElementById("linkClaim")
	link.innerHTML = "";
	linktxt.innerHTML = "";
}

function filesPlural(txt){
	var count = 0;
	for(var i = 0; i < txt.length; i++) {
		if (txt[i] == '\n'){
			count++
		}
		if (count >=2) {
			return true
		}
	}
	return false
}

function FinalRecord(){
	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			try {
				finalRecordObj = JSON.parse(txt);
			}
			catch(err) {
				var element=document.getElementById("toRecordEmail");
				var para = document.createElement("p")
				para.innerText = xmlhttp.responseText;
				element.appendChild(para)
				return
			}
			FinalRecordOutput()
		}
	}
	var job = document.getElementById("finalArtworkInput").value;
	xmlhttp.open("GET","/record?j="+job,true);
	xmlhttp.send();

}

function FinalRecordOutput(){
	if (finalartworkObj == null){
		return
	}
	ClearContent("toRecordEmail")

	var element=document.getElementById("toRecordEmail");
	var lineBreak = document.createElement("br");


	var emailHeader = document.createElement("p");
	emailHeader.innerText = "Files For Record WMT Canada / " + finalRecordObj.PHQ;
	emailHeader.setAttribute("id", "EmailHeadertxt");
	emailHeader.setAttribute("contenteditable", "true")
	element.appendChild(emailHeader)
	element.appendChild(lineBreak)
	element.appendChild(lineBreak.cloneNode())

	var emailBoby = document.createElement("p");

	var outPutTxt = finalartworkObj.Body;
	if (outPutTxt.length === 0){
		emailBoby.innerText = finalartworkObj.Error;
		element.appendChild(emailBoby);
		return
	}

	emailBoby.innerText = helloFaris + "\n" +outPutTxt + "\n"+ finalRecordObj.TxtBodies[0].TxtBody;
	emailBoby.setAttribute("id", "EmailBody");
	emailBoby.setAttribute("contenteditable", "true");
	element.appendChild(emailBoby);

	element.appendChild(lineBreak);

	var btn = document.createElement("button");
	btn.innerText = "autoEmail"
	btn.setAttribute('onclick', "autoEmail()");
	element.appendChild(btn);

}

var helloFaris = `
Hi Faris,

We’ve released final files to supplier and attached files are for your record:
`

function autoEmail(){
	var ElementEmailHeader=document.getElementById("EmailHeadertxt");
	var ElementEmailContent=document.getElementById("EmailBody");

	var EmailHeader = ElementEmailHeader.innerText;
	var EmailContent = ElementEmailContent.innerText;

	var URL_EmailHeaer = encodeURIComponent(EmailHeader);
	var URL_EmailContent = encodeURIComponent(EmailContent);

	var params = "title="+ URL_EmailHeaer + "&" + "content=" + URL_EmailContent;

	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			try {
				// finalRecordObj = JSON.parse(txt);
				//console.log(txt)
			}
			catch(err) {
				// var element=document.getElementById("toRecordEmail");
				// var para = document.createElement("p")
				// para.innerText = xmlhttp.responseText;
				// element.appendChild(para)
				return
			}
			// FinalRecordOutput()
			//console.log("test")
		}
	}
	// var job = document.getElementById("finalArtworkInput").value;
	xmlhttp.open("POST","/autoemail",true);
	xmlhttp.setRequestHeader('Content-type', 'application/x-www-form-urlencoded');
	xmlhttp.send(params);
}

function autoEmailButton(){
	ClearContent("toRecordEmail")
	var element=document.getElementById("toRecordEmail");
	var btn = document.createElement("button");
	btn.innerText = "autoRecord";
	btn.setAttribute('onclick', "FinalRecord()");
	element.appendChild(btn);
}