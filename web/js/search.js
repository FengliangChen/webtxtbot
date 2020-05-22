function supplierCodeUpdate(){
	var supplierCode = document.getElementById("supplierCode").value
	var supplierValue = document.getElementById("supplierValue").value
	if (supplierCode.length == 0 && supplierValue.length == 0) {
		alert("注意填写！！！")
		return
	}
	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			displayMessage(txt)
		}
	}

	xmlhttp.open("GET","/updatesupplier?c="+supplierCode+"&"+"v="+supplierValue,true);
	xmlhttp.send();
}

function supplierSearch(supplierCode){
	if (supplierCode.length < 2){
		document.getElementById("displaySearchResult").innerHTML = "";
		document.getElementById("returnMessage").innerHTML = "";
		document.getElementById("displaySearchError").innerHTML = "";
		return
	}

	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText;

			document.getElementById("displaySearchResult").innerHTML = "";

			if (txt.length == 0 ){ return;};

			try {
				var searchResultObj = JSON.parse(txt);
			}
			catch(err) {
				document.getElementById("displaySearchError").innerText = txt;

			}
			CreateForm(searchResultObj);
		}
	}

	xmlhttp.open("GET","/suppliersearch?c="+supplierCode,true);
	xmlhttp.send();
}

function displayMessage(txt){
	document.getElementById("returnMessage").innerText = txt;
}

function CreateForm(resultObj){
	var table="<tr><th>客号</th><th>客户</th></tr>";
	for (i = 0; i < resultObj.length; i++ ){
		for(x in resultObj[i]){
			table += "<tr><td>" + x + "</td><td>" + resultObj[i][x] +"</td></tr>";
		}
    	
	}
	document.getElementById("displaySearchResult").innerHTML = table;
}
