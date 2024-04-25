interface ChatAPISidebarResponse {
  Name: string;
  Id: number;
}
type Chatroom = ChatAPISidebarResponse;
type ChatMessage = {
  Name: string;
  Message: string;
  Id: number;
};

enum SideBarTab {
  chatrooms = 0,
  invites,
  users,
}

let currentChatroom: Chatroom;
let lastMessageId: number;

const openTabEl: HTMLDivElement = document.querySelector("#open-tab");
const messageBoxEl: HTMLDivElement = document.querySelector("#message-box");

const appendChatMessage = (message: ChatMessage) => {
  const messageDiv = document.createElement("div");
  const messageP = document.createElement("p");
  messageP.innerText = `${message.Name}: ${message.Message}`;
  messageDiv.appendChild(messageP);
  messageBoxEl.appendChild(messageDiv);
};

const populateChatroom = (messages: ChatMessage[]) => {
  while (messageBoxEl.firstChild) {
    messageBoxEl.removeChild(messageBoxEl.firstChild);
  }
  messages.forEach(appendChatMessage);
};

const getChatMessages = async (chatroom: Chatroom) => {
  const messages: ChatMessage[] = await fetch("/messages", {
    method: "POST",
    body: JSON.stringify({ Id: chatroom.Id, LastMessage: -1 }),
  })
    .then((response) => response.json() as Promise<ChatMessage[] | null>)
    .catch((err) => {
      console.error(err);
      return null;
    });
  if (messages != null) {
    populateChatroom(messages);
    lastMessageId = messages[messages.length - 1].Id;
  } else {
    populateChatroom([
      { Name: "", Message: "No messages in this chatroom", Id: -1 },
    ]);
    lastMessageId = -1;
  }
  (document.querySelector("#chatroom-name") as HTMLParagraphElement).innerText =
    chatroom.Name;
  currentChatroom = chatroom;
};

const checkNewMessages = async () => {
  const messages: ChatMessage[] = await fetch("/messages", {
    method: "POST",
    body: JSON.stringify({
      Id: currentChatroom.Id,
      LastMessage: lastMessageId,
    }),
  })
    .then((response) => response.json() as Promise<ChatMessage[]>)
    .catch((err) => {
      console.error(err);
      return null;
    });
  if (messages != null && messages.length != 0) {
    if (lastMessageId === -1) {
      populateChatroom(messages);
    } else {
      for (let message of messages) {
        appendChatMessage(message);
      }
    }
    lastMessageId = messages[messages.length - 1].Id;
  }
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
  if (tabIndex === SideBarTab.chatrooms) {
    (<HTMLButtonElement>(
      document.querySelector("#new-group-btn")
    )).style.display = "inline";
  } else {
    (<HTMLButtonElement>(
      document.querySelector("#new-group-btn")
    )).style.display = "none";
  }
};

const getChatrooms = async () => {
  setSidebarLoading();
  changeTabCss(SideBarTab.chatrooms);
  const chatrooms: ChatAPISidebarResponse[] = await fetch("/chatrooms")
    .then((response) => response.json() as Promise<ChatAPISidebarResponse[]>)
    .catch((err) => {
      console.error(err);
      return null;
    });
  if (chatrooms === null) {
    openTabEl.removeChild(openTabEl.firstChild);
    const p = document.createElement("p");
    p.innerText = "You are not in any chatrooms";
    openTabEl.appendChild(p);
  } else {
    openTabEl.removeChild(openTabEl.firstChild);
    chatrooms.forEach((chatroom) => {
      const chatroomEl = document.createElement("div");
      chatroomEl.classList.add("chatroom-preview", "preview");
      const nameEl = document.createElement("p");
      nameEl.innerText = chatroom.Name;
      const idEl = document.createElement("span");
      idEl.innerText = ` #${chatroom.Id}`;
      chatroomEl.addEventListener("click", () => {
        getChatMessages(chatroom).catch((err) => console.error(err));
      });

      nameEl.appendChild(idEl);
      chatroomEl.appendChild(nameEl);
      openTabEl.appendChild(chatroomEl);
    });
    if (currentChatroom == undefined) {
      getChatMessages(chatrooms[0]).catch((err) => console.error(err));
    }
  }
};

const getInvites = async () => {
  setSidebarLoading();
  changeTabCss(SideBarTab.invites);
  const invites: ChatAPISidebarResponse[] = await fetch("/invites")
    .then((response) => response.json() as Promise<ChatAPISidebarResponse[]>)
    .catch((err) => {
      console.error(err);
      return null;
    });
  if (invites === null) {
    openTabEl.removeChild(openTabEl.firstChild);
    const p = document.createElement("p");
    p.innerText = "You have no invites";
    openTabEl.appendChild(p);
  } else {
    openTabEl.removeChild(openTabEl.firstChild);
    invites.forEach((invite) => {
      const inviteEl = document.createElement("div");
      inviteEl.classList.add("invite-preview", "preview");
      const nameEl = document.createElement("p");
      nameEl.innerText = invite.Name;
      nameEl.addEventListener("click", () => {
        if (confirm(`accept invite to ${invite.Name}`)) {
          fetch("/acceptInvite", {
            method: "POST",
            body: JSON.stringify({ InviteId: invite.Id }),
          })
            .then((response) => {
              if (response.status === 200)
                getInvites().catch((err) => console.error(err));
              else console.log("failed to accept invite");
            })
            .catch((err) => console.error(err));
        }
      });
      inviteEl.appendChild(nameEl);
      openTabEl.appendChild(inviteEl);
    });
  }
};

const getUsers = async () => {
  setSidebarLoading();
  changeTabCss(SideBarTab.users);
  const users: null | ChatAPISidebarResponse[] = await fetch("/users")
    .then((response) => response.json() as Promise<ChatAPISidebarResponse[]>)
    .catch((err) => {
      console.error(err);
      return null;
    });
  if (users === null) return;
  if (users.length > 0) {
    openTabEl.removeChild(openTabEl.firstChild);
    users.forEach((user) => {
      const userEl = document.createElement("div");
      userEl.classList.add("user-preview", "preview");
      const nameEl = document.createElement("p");
      nameEl.innerText = user.Name;
      userEl.addEventListener("click", () => {
        if (confirm(`invite ${user.Name} to ${currentChatroom.Name}`)) {
          console.log("send invite");
          fetch("/sendInvite", {
            method: "POST",
            body: JSON.stringify({
              ChatId: currentChatroom.Id,
              UserId: user.Id,
            }),
          })
            .then((response) => {
              if (response.status !== 200) console.log("failed to send invite");
            })
            .catch((err) => console.error(err));
        }
      });
      userEl.appendChild(nameEl);
      openTabEl.appendChild(userEl);
    });
  }
};

const sendMessage = async (e: Event) => {
  e.preventDefault();
  const message = (
    document.querySelector('textarea[name="Message"]') as HTMLTextAreaElement
  ).value;
  const response: Response = await fetch("/send", {
    method: "POST",
    body: JSON.stringify({ ChatId: currentChatroom.Id, Message: message }),
  }).catch((err) => {
    console.error(err);
    return null;
  });
  if (response !== null && response.status === 200) {
    getChatMessages(currentChatroom).catch((err) => console.error(err));
    (
      document.querySelector('textarea[name="Message"]') as HTMLTextAreaElement
    ).value = "";
  }
};

const createNewChat = async (e: Event) => {
  e.preventDefault();
  const name = (<HTMLInputElement>(
    document.querySelector('input[name="new-chatroom-name"]')
  )).value;
  const response = await fetch("/createChatroom", {
    method: "POST",
    body: JSON.stringify({ Name: name }),
  });
};

const openNewChatModal = (e: Event) => {
  e.preventDefault();
  (<HTMLDivElement>(
    document.querySelector("#new-chatroom-modal")
  )).style.display = "flex";
};
const closeNewChatModal = (e: Event) => {
  e.preventDefault();
  (<HTMLDivElement>(
    document.querySelector("#new-chatroom-modal")
  )).style.display = "none";
};

document
  .querySelector("#chatroomsTabBtn")
  .addEventListener("click", getChatrooms);
document.querySelector("#invitesTabBtn").addEventListener("click", getInvites);
document.querySelector("#usersTabBtn").addEventListener("click", getUsers);
document.querySelector("#message-form").addEventListener("submit", sendMessage);
document
  .querySelector("#new-group-btn")
  .addEventListener("click", openNewChatModal);
document
  .querySelector("#chatroom-modal-cancel")
  .addEventListener("click", closeNewChatModal);
document
  .querySelector("#submit-new-chatroom")
  .addEventListener("click", (e) => {
    (async (e) => {
      await createNewChat(e).catch((err) => console.error(err));
      closeNewChatModal(e);
      await getChatrooms().catch((err) => console.error(err));
    })(e);
  });
getChatrooms().catch((err) => console.error(err));

setInterval(async () => {
  if (currentChatroom != null) {
    await checkNewMessages();
  }
}, 10000);
