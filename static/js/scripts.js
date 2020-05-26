
function storeLSAT(id) {
	var preimage = document.getElementById("preimage").value;
	
	if (preimage == "") {
		alert("Empty preimage field");
		return;
	}

	var macaroon = document.getElementById("macaroon").value;
	var lsat = macaroon+":"+preimage;
	var loc = 'snell/article/'+ id.toString();

	localStorage.setItem(loc, lsat);
	
	var url = '/article/view/' + id.toString();
	
	xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function () {	
		if (xhr.readyState === 4) {
			document.body.innerHTML = '';
      			document.write(xhr.response);
   		}
	};
	xhr.open("GET", url, true);
	xhr.setRequestHeader("Authorization","LSAT "+lsat);
	xhr.send();
}

function getAuth(id) {	
	var loc = 'snell/article/'+ id.toString();
	var lsat = localStorage.getItem(loc);
	
	var url = '/article/view/' + id.toString();
	
	xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function () {	
		if (xhr.readyState === 4) {
      			document.write(xhr.response);
   		}
	};
	xhr.open("GET", url, true);
	xhr.setRequestHeader("Authorization","LSAT "+lsat);
	xhr.send();
}
