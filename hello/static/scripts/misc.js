$(document).ready(function(){
	$('[data-toggle="tooltip"]').tooltip();   
});

function jobInProgress(uid, callback) {
	var xhr = new XMLHttpRequest();
	xhr.open('GET', 'jobInfo?uid=' + uid);
	xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
	xhr.onload = function() {
		console.log('Job in progress : ' + xhr.responseText);
		var resp = JSON.parse(xhr.responseText);
		callback( resp.Source_id != '')
	};
	xhr.send();

}

function manageSections(jobInProgress) {
	if(!jobInProgress){
		$("#selectSection").collapse('show');

		var sourceName = document.getElementById("sourceName").innerHTML
		var destName = document.getElementById("destName").innerHTML

		if( sourceName != "" && destName != "") {
			$("#transferButtonSection").collapse('show');
		} else {
			$("#transferButtonSection").collapse('hide');
		}

	} else {
		$("#jobSection").collapse('show');
		if (window.Worker){
			var progressUpdater = new Worker("scripts/progressUpdater.js")


			progressUpdater.onmessage = function(e) {
				console.log(e.data.percentage)


				$("#jobProgressBar").css('width', e.data.percentage + '%');
				$('#jobProgressBar').html(Math.floor(e.data.percentage) + '%');

				if (e.data.percentage > 0) {
					$('#initializingTransfer').html("Email threads processed: " + e.data.processed + "/" + e.data.total);
				}

			};

			var uidMessage = { uid: uid };
			progressUpdater.postMessage(uidMessage)

		}
	}
}
