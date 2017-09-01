
function setElements(isLoggedIn){
  var auth2 = gapi.auth2.getAuthInstance();
  var googleUser= auth2.currentUser.get();
  if(isLoggedIn){
      document.getElementById('gSignInButton').style.display = 'none';
      document.getElementById('logout').style.display = 'block';
      $("#selectSection").collapse('show');

  } else {
      $("#selectSection").collapse('hide');
      document.getElementById('logout').style.display = 'none';
      document.getElementById('gSignInButton').style.display = 'block';
  }
}

function signOut() {
  var auth2 = gapi.auth2.getAuthInstance();
  auth2.signOut().then(function () {
    setElements(false);
    console.log('User signed out.');
  });
  
}

function sendTokentoDB(id_token){
  var auth2 = gapi.auth2.getAuthInstance();
  var googleUser= auth2.currentUser.get();
  var profile = googleUser.getBasicProfile();
  
  var xhr = new XMLHttpRequest();
  xhr.open('POST', 'signIn');
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onload = function() {
    setElements(true);
  };

  var data = {
      Uid : id_token
    , Name : profile.getName()
  }

  xhr.send(JSON.stringify(data));

  console.log("Sent: " + JSON.stringify(data) + " to database");

}

//after sign in, verify token
function onSignIn(googleUser) {
  document.getElementById('gSignInButton').style.display = 'none';
  var id_token = googleUser.getAuthResponse().id_token;

  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=' + id_token);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.onload = function() {
    console.log('Signed in as: ' + xhr.responseText);
    var resp = JSON.parse(xhr.responseText);

    //token is valid, send to back end
    if (resp.aud === "65587295914-kbl4e2chuddg9ml7d72f6opqhddl62fv.apps.googleusercontent.com") {
      sendTokentoDB(resp.sub);
    } else {
      signOut();
    }
  };
  xhr.send();
  

  var profile = googleUser.getBasicProfile();
  console.log('ID: ' + profile.getId()); // Do not send to your backend! Use an ID token instead.
  console.log('Name: ' + profile.getName());
  console.log('Image URL: ' + profile.getImageUrl());
  console.log('Email: ' + profile.getEmail()); // This is null if the 'email' scope is not present.
  console.log('id_token: ' + id_token);
  
}

function askPermissions(accountType) {

  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'askPermissions?type=' + accountType);
  xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
  xhr.onload = function() {
    console.log('Got ' + xhr.responseText + ' from server');
  };
  xhr.send();

  console.log("Sent: askPermissions?type=" + accountType + " to server");

}

(function() {
    var cors_api_host = 'cors-anywhere.herokuapp.com';
    var cors_api_url = 'https://' + cors_api_host + '/';
    var slice = [].slice;
    var origin = window.location.protocol + '//' + window.location.host;
    var open = XMLHttpRequest.prototype.open;
    XMLHttpRequest.prototype.open = function() {
        var args = slice.call(arguments);
        var targetOrigin = /^https?:\/\/([^\/]+)/i.exec(args[1]);
        if (targetOrigin && targetOrigin[0].toLowerCase() !== origin &&
            targetOrigin[1] !== cors_api_host) {
            args[1] = cors_api_url + args[1];
        }
        return open.apply(this, args);
    };
})();