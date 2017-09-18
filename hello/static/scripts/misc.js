$(document).ready(function(){
    $('[data-toggle="tooltip"]').tooltip();   
});

function jobInProgress(uid, callback) {
  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'jobInProgress?uid=' + uid);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.onload = function() {
    console.log('Job in progress : ' + xhr.responseText);
    var resp = JSON.parse(xhr.responseText);
    callback( resp.InProgress == 'true')
  };
  xhr.send();
  
}
