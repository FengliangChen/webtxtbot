
var seasonNameobj
var seasonKeys
var buyerList

if (seasonNameobj == null) {
	GetSeasonMapJson()
}

if (buyerList == null) {
	GetBuyerListJson()
}


function SeasonSelect(value){
		// document.getElementById('seasonDatalist').innerHTML = '';

		// var l=value.length;
		// var n = seasonKeys.length;
		// for (var i = 0; i<n; i++) {
		// 	if(((seasonKeys[i].toUpperCase()).indexOf(value.toUpperCase()))>-1){
		// 		var node = document.createElement("option");
		// 		var val = document.createTextNode(seasonKeys[i]);
		// 		node.appendChild(val);
		// 		document.getElementById('seasonDatalist').appendChild(node);
		// 	}
		// }


		// for (var i = 0; i < seasonKeys.length; i++ ){
		// 	var node = document.createElement("option");
		// 	var val = document.createTextNode(seasonKeys[i]);
		// 	node.appendChild(val);
		// 	document.getElementById('seasonDatalist').appendChild(node);

		// }

		UpdateProgramInput(value.toUpperCase())
		// GetSeasonMapJson()
	}

function BuyerSelect(value){
		document.getElementById('buyerlist').innerHTML = '';

		// var l=value.length;
		// var n = seasonKeys.length;
		// for (var i = 0; i<n; i++) {
		// 	if(((seasonKeys[i].toUpperCase()).indexOf(value.toUpperCase()))>-1){
		// 		var node = document.createElement("option");
		// 		var val = document.createTextNode(seasonKeys[i]);
		// 		node.appendChild(val);
		// 		document.getElementById('seasonDatalist').appendChild(node);
		// 	}
		// }


		for (var i = 0; i < buyerList.length; i++ ){
			var node = document.createElement("option");
			var val = document.createTextNode(buyerList[i]);
			node.appendChild(val);
			document.getElementById('buyerlist').appendChild(node);

		}
	}

function UpdateProgramInput(seasonCode){
	var element=document.getElementById("programInput");
	element.value = seasonNameobj[seasonCode];
}

function GetSeasonMapJson(){
	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			try{
				seasonNameobj = JSON.parse(txt)
			}
			catch(err){
				alert(txt)
			}
			seasonKeys = Object.keys(seasonNameobj)

			// for (var i = 0; i < seasonKeys.length; i++ ){
			// 	console.log(seasonKeys[i])

			// }
		}
	}
	xmlhttp.open("GET","/job/seasonmap",true);
	xmlhttp.send();
}

function GetBuyerListJson(){
	var xmlhttp;
	xmlhttp=new XMLHttpRequest();
	xmlhttp.onreadystatechange=function(){
		if (xmlhttp.readyState==4 && xmlhttp.status==200)
		{
			var txt = xmlhttp.responseText
			try{
				var buyerobj = JSON.parse(txt)
			}
			catch(err){
				alert(txt)
			}

			buyerList = buyerobj

			// for (var i = 0; i < seasonKeys.length; i++ ){
			// 	console.log(seasonKeys[i])

			// }
		}
	}
	xmlhttp.open("GET","/job/buyerlist",true);
	xmlhttp.send();
}