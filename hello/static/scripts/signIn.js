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
    console.log('User signed out.');
  });
  setElements(false);
}