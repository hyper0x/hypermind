$(function () {
    
  $('#cv_auth_code_submit').bind('click', function() {
    var auth_code = $('#cv_auth_code').val();
    //alert("Input: " + auth_code);
    $.post("/get-cv", "auth_code=" + auth_code,
        function (data, textStatus){
            if ((data.indexOf("FAIL:") == 0) || (data.indexOf("ERROR:") == 0)) {
                alert(data);  
            } else {
                openWindow = window.open("", "", "height=600, width=800,top=50,left=50,toolbar=no,menubar=no,scrollbars=auto,resizeable=no,location=no,status=no");  
                openWindow.document.write(data)
                openWindow.document.close(); 
            }
        }, "text");
    });

})