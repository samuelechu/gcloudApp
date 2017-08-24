
function setElements(isLoggedIn){
  if(isLoggedIn){
      document.getElementById('gSignInButton').style.display = 'none'
      document.getElementById('logout').style.display = 'block'

  } else {
      document.getElementById('gSignInButton').style.display = 'block'
      document.getElementById('logout').style.display = 'none'
  }
}

function signOut() {
  var auth2 = gapi.auth2.getAuthInstance();
  auth2.signOut().then(function () {
    $('.collapse').collapse("hide")
    console.log('User signed out.');
  });
  setElements(false);
}

function sendTokentoDB(googleUser, id_token){

  var profile = googleUser.getBasicProfile();
  
  var xhr = new XMLHttpRequest();
  xhr.open('POST', 'signIn');
  xhr.setRequestHeader("Content-Type", "application/json");
  xhr.onload = function() {
    $('.collapse').collapse("show")
  };

  var data = {
      Uid : id_token
    , Name : profile.getName()
  }

  xhr.send(JSON.stringify(data));

  console.log("Sent: " + JSON.stringify(data) + " to database");

}

function onSignIn(googleUser) {
  var id_token = googleUser.getAuthResponse().id_token;

  var xhr = new XMLHttpRequest();
  xhr.open('GET', 'https://www.googleapis.com/oauth2/v3/tokeninfo?id_token=' + id_token);
  xhr.setRequestHeader('Content-Type', 'application/x-www-form-urlencoded');
  xhr.onload = function() {
    console.log('Signed in as: ' + xhr.responseText);
    var resp = JSON.parse(xhr.responseText);

    //token is valid, send to back end
    if (resp.aud === "65587295914-kbl4e2chuddg9ml7d72f6opqhddl62fv.apps.googleusercontent.com") {
      sendTokentoDB(googleUser, resp.sub);
    }
  };
  xhr.send();
  

  var profile = googleUser.getBasicProfile();
  console.log('ID: ' + profile.getId()); // Do not send to your backend! Use an ID token instead.
  console.log('Name: ' + profile.getName());
  console.log('Image URL: ' + profile.getImageUrl());
  console.log('Email: ' + profile.getEmail()); // This is null if the 'email' scope is not present.
  console.log('id_token: ' + id_token);
  setElements(true);
}