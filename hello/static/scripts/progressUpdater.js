function timedCount()
{
i=i+1;
postMessage(i);                   //posts a message back to the HTML page.
setTimeout("timedCount()",500);
}

timedCount();




function updateProgress(){
	var xhr = new XMLHttpRequest();
	xhr.open('GET', 'progress?uid=' + uid);
	xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
	xhr.onload = function() {
	console.log('Job in progress : ' + xhr.responseText);
	var resp = JSON.parse(xhr.responseText);
	callback( resp.InProgress == 'true')
	};
	xhr.send();
  
}