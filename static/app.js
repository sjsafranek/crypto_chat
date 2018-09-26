
function randomHash() {
    return CryptoJS.MD5(
                new Date().toISOString()
            ).toString();
}

var app = {
    "chatrooms": {}
}

$('.newChatRoom').on('click', function(){
    swal({
        title: 'Join Chatroom',
        html: '<input id="chatroom_name" class="swal2-input" placeholder="chatroom">' +
            '<input id="chatroom_passphrase" class="swal2-input" placeholder="passphrase">',
        showCancelButton: true,
        focusConfirm: false,
        preConfirm: () => {
            return {
                chatroom: document.getElementById('chatroom_name').value,
                passphrase: document.getElementById('chatroom_passphrase').value
            }
        }
    }).then(function(value){
        if (value) {
            var result = value.value;

            if (app.chatrooms[result.chatroom]) {
                M.toast({html: 'Chatroom already open', displayLength: 3000});
                return;
            }

            app.chatrooms[result.chatroom] = new ChatRoom(result.chatroom, result.passphrase, function() {
                delete app.chatrooms[result.chatroom];
                removeTab(result.chatroom);
            });

            addTab(result.chatroom).append(
                app.chatrooms[result.chatroom].$el
            );
        }
    });
});




function buildTabs() {
    var $tabs = $('.tabs');
    var instance = M.Tabs.getInstance($tabs);
    instance && instance.destroy();
    $tabs.children().length && M.Tabs.init($tabs, {});
}

function removeTab(tab_name) {
    if (!tab_name) {
        M.toast({html: 'Must supply a tab name', displayLength: 3000});
        return;
    }
    var tab_id = CryptoJS.MD5(tab_name).toString();
    $('#'+tab_id).remove();
    buildTabs();
}

function addTab(tab_name) {
    if (!tab_name) {
        M.toast({html: 'Must supply a tab name', displayLength: 3000});
        return;
    }

    var tab_id = CryptoJS.MD5(tab_name).toString();

    // make singleton
    removeTab(tab_name);

    var $tabLabel = $('<li>').addClass('tab').append(
        $('<a>', {
            href: '#'+tab_id
        }).text(tab_name)
    );

    $('#chatrooms .tabs').append( $tabLabel );

    var $content = $('<div>', {id: tab_id})
                        .addClass('col s12')
                        .on('remove', function(){
                            $tabLabel.remove();
                        });

    $('#chatrooms').append(
        $content
    );

    // rebuild tabs
    buildTabs();

    // select new tab
    var $tabs = $('.tabs');
    var instance = M.Tabs.getInstance($tabs);
    instance.select(tab_id);

    // return tab content
    return $content;
}
