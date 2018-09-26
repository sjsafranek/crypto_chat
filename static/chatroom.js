
// strips out html...
function sanitizeText(text) {
    return $('<p>').html(text).text();
}

function gravatarURL(email) {
    return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
}

var ChatRoom = function(chatroom, passphrase, onClose){
    this.ws = null;
    this.data = {
        email: user.email,
        username: user.username,
        passphrase: passphrase,
        chatroom: chatroom
    }
    this.onClose = onClose;
    this.open();
    this.buildUi();
}

ChatRoom.prototype.open = function() {
    var self = this;

    this.ws = new WebSocket('ws://' + window.location.host + '/ws/' + this.data.chatroom);

    this.ws.addEventListener('message', function(e) {
        try {
            var result = JSON.parse(e.data);
            var decrypted = CryptoJS.AES.decrypt(result.data, self.data.passphrase);
            var msg = decrypted.toString(CryptoJS.enc.Utf8);
            msg = JSON.parse(msg);
            // write message to page
            self.chatMessages.append(
                '<div class="chip">'
                    + '<img src="' + gravatarURL(msg.email) + '">' // Avatar
                    + msg.username
                + '</div>'
                // Parse emojis
                + emojione.toImage(msg.message) + '<br/>'
            );
            // Auto scroll to the bottom
            self.chatMessages.scrollTop = self.chatMessages.scrollHeight;
        }
        catch(err) {
            console.log(err)
            M.toast({html: 'Unable to parse incomming message', displayLength: 3000});
            return;
        }
    });

    this.ws.onopen = function(e) {
        M.toast({html: 'Connection open', displayLength: 2000});
    }

    this.ws.onclose = function(e) {
        M.toast({html: 'Connection closed', displayLength: 5000});
    }
}

ChatRoom.prototype.destroy = function() {
    this.$el.remove();
}

ChatRoom.prototype.close = function() {
    this.ws && this.ws.close();
    this.destroy();
    this.onClose && this.onClose();
}

ChatRoom.prototype.send = function() {
    var newMsg = this.inputMessage.val();
    if (newMsg != '') {

        var data = JSON.stringify({
            email: this.data.email,
            username: this.data.username,
            message: sanitizeText(newMsg)
        });

        var encrypted = CryptoJS.AES.encrypt(data, this.data.passphrase);
        var payload = encrypted.toString();

        this.ws.send(
            JSON.stringify({"data": payload})
        );

        this.inputMessage.val('');
    }

}

ChatRoom.prototype.buildUi = function() {
    var self = this;

    this.chatMessages = $('<div>').addClass('card-content chat-messages');
    this.inputMessage = $('<input>', {type: 'text', id:'message'}).addClass('validate');

    this.$el = $('<div>').addClass('row').append(
        // $('<div>').addClass('col s12 m6').append(
        $('<div>').addClass('col s12').append(
            $('<div>').addClass('card').append(
                $('<div>').addClass('card-title').append(
                    $('<span>').append(this.data.chatroom),
                    $('<button>')
                            .addClass('right waves-effect waves-light btn red btn-small')
                            .on('click', function(){
                                self.close();
                            })
                            .append(
                                $('<i>').addClass('material-icons').append('close')
                            )
                ),
                this.chatMessages,
                $('<div>').addClass('card-action').append(
                    $('<div>').addClass('input-field input-group').append(
                        this.inputMessage,
                        $('<label>', {for:'message'}).append('Message'),
                        $('<span>').addClass('suffix').append(
                            $('<a>').addClass('waves-effect waves-light btn btn-floating')
                                .on('click', function() {
                                    self.send();
                                })
                                .append(
                                    $('<i>').addClass('material-icons').append('chat')
                                )
                        )
                    )
                )
            )
        )
    );
}
