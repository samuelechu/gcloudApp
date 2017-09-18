self.addEventListener("message", function(e) {
	updateProgress(e.data.uid)
    // the passed-in data is available via e.data
}, false);

function updateProgress(uid){

	var xhr = new XMLHttpRequest();
	xhr.open('GET', '../../jobInfo?uid=' + uid);
	xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  	xhr.onload = function() {
    console.log('Job in progress : ' + xhr.responseText);
    var resp = JSON.parse(xhr.responseText);
   
    var percentage = resp.Processed_threads * 100 / resp.Total_threads
    console.log('percentage: ' + percentage)
    var percentageMessage = { percentage: percentage};

    postMessage(percentageMessage)
    if(percentage < 1){
		setTimeout("updateProgress()",5000);
	}
  };
  xhr.send();
  
}