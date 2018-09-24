
// strips out html...
function sanitizeText(text) {
    return $('<p>').html(text).text();
}

var chatApp = new Vue({
    el: '#app',

    data: {
        ws: null, // Our websocket
        newMsg: '', // Holds new messages to be sent to the server
        chatContent: '', // A running list of chat messages displayed on the screen
        email: null, // Email address used for grabbing an avatar
        username: null, // Our username
        passphrase: null,
        chatroom: null,
        debug: true,
        joined: false // True if email and username have been filled in
    },

    created: function() {
        // var self = this;
        // this.ws = new WebSocket('ws://' + window.location.host + '/ws');
        // this.ws.addEventListener('message', function(e) {
        //     console.log(e.data);
        //
        //     try {
        //         var result = JSON.parse(e.data);
        //         var decrypted = CryptoJS.AES.decrypt(result.data, self.passphrase);
        //         var msg = decrypted.toString(CryptoJS.enc.Utf8);
        //         msg = JSON.parse(msg);
        //
        //         // write message to page
        //         self.chatContent += '<div class="chip">'
        //                 + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
        //                 + msg.username
        //             + '</div>'
        //             + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
        //
        //         var element = document.getElementById('chat-messages');
        //         element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
        //
        //     }
        //     catch(err) {
        //         Materialize.toast('Unable to parse incomming message', 3000);
        //         return;
        //     }
        //
        // });
        //
        // this.ws.onopen = function(e) {
        //     console.log(e);
        // }
    },

    methods: {
        send: function () {
            if (this.newMsg != '') {

                var data = JSON.stringify({
                    email: this.email,
                    username: this.username,
                    message: sanitizeText(this.newMsg)
                });

                var encrypted = CryptoJS.AES.encrypt(data, this.passphrase);
                var payload = encrypted.toString();

                this.ws.send(
                    JSON.stringify({"data": payload})
                );
                this.newMsg = ''; // Reset newMsg
            }
        },

        join: function () {
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            if (!this.chatroom) {
                Materialize.toast('You must choose a chat room', 2000);
                return
            }
            this.email = sanitizeText(this.email);
            this.username = sanitizeText(this.username);
            this.passphrase = sanitizeText(this.passphrase);
            this.chatroom = sanitizeText(this.chatroom);
            this.joined = true;

            this.connect();
        },

        connect: function(){
            var self = this;

            this.ws = new WebSocket('ws://' + window.location.host + '/ws/' + this.chatroom);

            this.ws.addEventListener('message', function(e) {
                console.log(e.data);
                try {
                    var result = JSON.parse(e.data);
                    var decrypted = CryptoJS.AES.decrypt(result.data, self.passphrase);
                    var msg = decrypted.toString(CryptoJS.enc.Utf8);
                    msg = JSON.parse(msg);
                    // write message to page
                    self.chatContent += '<div class="chip">'
                            + '<img src="' + self.gravatarURL(msg.email) + '">' // Avatar
                            + msg.username
                        + '</div>'
                        + emojione.toImage(msg.message) + '<br/>'; // Parse emojis
                    var element = document.getElementById('chat-messages');
                    element.scrollTop = element.scrollHeight; // Auto scroll to the bottom
                }
                catch(err) {
                    Materialize.toast('Unable to parse incomming message', 3000);
                    return;
                }
            });

            this.ws.onopen = function(e) {
                Materialize.toast('Connection open', 2000);
            }

            this.ws.onclose = function(e) {
                Materialize.toast('Connection closed', 15000);
            }

        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
        }
    }
});