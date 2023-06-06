interface ChatAPISidebarResponse {
  Name: string;
  Id: number;
}
type Chatroom = ChatAPISidebarResponse;
type ChatMessage = {
  Name: string;
  Message: string;
};

enum SideBarTab {
  chatrooms = 0,
  invites,
  users,
}

let currentChatroom: Chatroom;

const openTabEl: HTMLDivElement = document.querySelector("#open-tab");
const messageBoxEl: HTMLDivElement = document.querySelector("#message-box");

const populateChatroom = (messages: ChatMessage[]) => {
  while (messageBoxEl.firstChild) {
    messageBoxEl.removeChild(messageBoxEl.firstChild);
  }
  messages.forEach((message) => {
    const messageDiv = document.createElement("div");
    const messageP = document.createElement("p");
    messageP.innerText = `${message.Name}: ${message.Message}`;
    messageDiv.appendChild(messageP);
    messageBoxEl.appendChild(messageDiv);
  });
};

const getChatMessages = (chatroom: Chatroom) => {
  fetch("/messages", {
    method: "POST",
    body: JSON.stringify({ Id: chatroom.Id }),
  })
    .then((response) => response.json() as Promise<ChatMessage[] | null>)
    .then((messages) => {
      if (messages != null) populateChatroom(messages);
      else {
        console.log("no messages");
      }
    })
    .catch((err) => console.error(err));
};

const setSidebarLoading = () => {
  while (openTabEl.firstChild) {
    openTabEl.removeChild(openTabEl.firstChild);
  }
  const loading = document.createElement("p");
  loading.innerText = "loading...";
  openTabEl.appendChild(loading);
};

const changeTabCss = (tabIndex: SideBarTab) => {
  (<NodeListOf<HTMLParagraphElement>>(
    document.querySelectorAll(".tabs>p")
  )).forEach((tab: HTMLParagraphElement, i: number) => {
    if (i === tabIndex) {
      tab.classList.add("selected-tab");
    } else {
      tab.classList.remove("selected-tab");
    }
  });
};

const getChatrooms = () => {
  setSidebarLoading();
  changeTabCss(SideBarTab.chatrooms);
  fetch("/chatrooms")
    .then((response) => response.json() as Promise<ChatAPISidebarResponse[]>)
    .then((chatrooms) => {
      if (chatrooms.length > 0) {
        openTabEl.removeChild(openTabEl.firstChild);
        chatrooms.forEach((chatroom) => {
          const chatroomEl = document.createElement("div");
          chatroomEl.classList.add("chatroom-preview", "preview");
          const nameEl = document.createElement("p");
          nameEl.innerText = chatroom.Name;
          const idEl = document.createElement("span");
          idEl.innerText = ` #${chatroom.Id}`;
          chatroomEl.addEventListener("click", () => {
            getChatMessages(chatroom);
            (
              document.querySelector("#chatroom-name") as HTMLParagraphElement
            ).innerText = chatroom.Name;
            currentChatroom = chatroom;
          });
          nameEl.appendChild(idEl);
          chatroomEl.appendChild(nameEl);
          openTabEl.appendChild(chatroomEl);
        });
      } else {
        console.log("no chat rooms");
      }
    })
    .catch((err) => {
      console.error(err);
    });
};

const getInvites = () => {
  setSidebarLoading();
  changeTabCss(SideBarTab.invites);
  fetch("/invites")
    .then((response) => response.json() as Promise<ChatAPISidebarResponse[]>)
    .then((invites) => {
      if (invites.length > 0) {
        openTabEl.removeChild(openTabEl.firstChild);
        invites.forEach((invite) => {
          const inviteEl = document.createElement("div");
          inviteEl.classList.add("invite-preview", "preview");
          const nameEl = document.createElement("p");
          nameEl.innerText = invite.Name;
          nameEl.addEventListener("click", () => {
            if (confirm(`accept invite to ${invite.Name}`)) {
              console.log(invite.Id);
            }
          });
          inviteEl.appendChild(nameEl);
          openTabEl.appendChild(inviteEl);
        });
      } else {
        console.log("no chat rooms");
      }
    })
    .catch((err) => {
      console.error(err);
    });
};

const getUsers = () => {
  setSidebarLoading();
  changeTabCss(SideBarTab.users);
  fetch("/users")
    .then((response) => response.json() as Promise<ChatAPISidebarResponse[]>)
    .then((users) => {
      if (users.length > 0) {
        openTabEl.removeChild(openTabEl.firstChild);
        users.forEach((user) => {
          const userEl = document.createElement("div");
          userEl.classList.add("user-preview", "preview");
          const nameEl = document.createElement("p");
          nameEl.innerText = user.Name;
          userEl.addEventListener("click", () => {
            if (confirm(`invite ${user.Name} to open group`)) {
              console.log(user.Id);
            }
          });
          userEl.appendChild(nameEl);
          openTabEl.appendChild(userEl);
        });
      } else {
        console.log("no chat rooms");
      }
    })
    .catch((err) => {
      console.error(err);
    });
};

const sendMessage = (e: Event) => {
  e.preventDefault();
  const message = (
    document.querySelector('textarea[name="Message"]') as HTMLTextAreaElement
  ).value;
  console.log(currentChatroom, message);
  fetch("/send", {
    method: "POST",
    body: JSON.stringify({ ChatId: currentChatroom.Id, Message: message }),
  })
    .then((response) => response.json() as Promise<{ Ok: boolean }>)
    .then((response) => {
      if (response.Ok) {
        getChatMessages(currentChatroom);
        (
          document.querySelector(
            'textarea[name="Message"]'
          ) as HTMLTextAreaElement
        ).value = "";
      }
    })
    .catch((err) => console.error(err));
};

const handleLogout = () => {
  console.log("logout");
  fetch("/logout")
    .then((response) => response.json())
    .then((response) => console.log(response))
    .catch((err) => console.log(err));
};

document
  .querySelector("#chatroomsTabBtn")
  .addEventListener("click", getChatrooms);
document.querySelector("#invitesTabBtn").addEventListener("click", getInvites);
document.querySelector("#usersTabBtn").addEventListener("click", getUsers);
document.querySelector("#message-form").addEventListener("submit", sendMessage);
document.querySelector("#logout-btn").addEventListener("click", handleLogout);

getChatrooms();
