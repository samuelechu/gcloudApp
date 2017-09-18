
function updateProgress(uid){
	var auth2 = gapi.auth2.getAuthInstance();
	var googleUser= auth2.currentUser.get();
	var profile = googleUser.getBasicProfile()
	var uid = profile.getId()
	
	var xhr = new XMLHttpRequest();
	xhr.open('GET', 'jobInfo?uid=' + uid);
	xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  	xhr.onload = function() {
    console.log('Job in progress : ' + xhr.responseText);
    var resp = JSON.parse(xhr.responseText);
    var percentage = resp.Processed_threads * 100 / resp.Total_threads
    console.log('percentage: ' + percentage)

    postMessage(percentage)
    if(percentage < 1){
		setTimeout("updateProgress()",5000);
	}
  };
  xhr.send();
  
}

updateProgress();