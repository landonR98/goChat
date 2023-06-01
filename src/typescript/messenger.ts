let chatroomsList;
let invitesList;
let usersList;

window.addEventListener("DOMContentLoaded", () => {
  const sideBarDisplay: HTMLDivElement =
    document.querySelector("#sidebar-display");

  const setSidebarLoading = () => {
    while (sideBarDisplay.firstChild) {
      sideBarDisplay.removeChild(sideBarDisplay.firstChild);
    }
    const loading = document.createElement("p");
    loading.innerText = "loading...";
    sideBarDisplay.appendChild(loading);
  };

  document
    .querySelector("#chatroomsTabBtn")
    .addEventListener("onclick", (event: Event) => {
      setSidebarLoading();
      fetch("/chatrooms")
        .then((response) => response.json())
        .then((json) => console.log(json))
        .catch((err) => {
          console.error(err);
        });
    });
  document
    .querySelector("#invitesTabBtn")
    .addEventListener("onclick", (event: Event) => {
      setSidebarLoading();
      fetch("/invites")
        .then((response) => response.json())
        .then((json) => console.log(json))
        .catch((err) => {
          console.error(err);
        });
    });
  document
    .querySelector("#usersTabBtn")
    .addEventListener("onclick", (event: Event) => {
      setSidebarLoading();
      fetch("/users")
        .then((response) => response.json())
        .then((json) => console.log(json))
        .catch((err) => {
          console.error(err);
        });
    });
});
