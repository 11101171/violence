// var socket;
// var resp = ''
// $(document).ready(function () {
//     // Create a socket
//     socket = new WebSocket('ws://' + window.location.host + '/admin/ssh/json?agentId=' + $('#agentId').val());
//     // Message received on the socket

//     socket.onmessage = function (event) {
//         var data = JSON.parse(event.data);
//         console.log(data);
//         switch (data.Type) {
//         case 0: // JOIN
//             // if (data.User == $('#uname').text()) {
//             //     $("#chatbox li").first().before("<li>You joined the chat room.</li>");
//             // } else {
//             //     $("#chatbox li").first().before("<li>" + data.User + " joined the chat room.</li>");
//             // }
//             receivedContent(data);
//             break;
//         case 1: // LEAVE
//             // $("#chatbox li").first().before("<li>" + data.User + " left the chat room.</li>");
//             receivedContent(data);
//             break;
//         case 2: // MESSAGE
//             // $("#chatbox li").first().before("<li><b>" + data.User + "</b>: " + data.Content + "</li>");
//             receivedContent(data);
//             break;
//         }
//     };

// });

//    // Send messages.
//     var postContent = function (content) {
//         socket.send(content);
//         var result = ''
//         var count = 1;
//         while(result == '' && count<10000) {
//             result = resp;
//             count++;
//         }
//         resp = '';
//         return result;
//     }

//     var receivedContent = function(data){
//         console.log('receivedContent'+ data.Type + data.User + data.Content )
//         resp = data.Content;
//     }