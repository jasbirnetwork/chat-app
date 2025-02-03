document.addEventListener("DOMContentLoaded", function () {
  const chatBox = document.getElementById("chat-box");
  const chatForm = document.getElementById("chat-form");
  const messageInput = document.getElementById("message-input");
  const userListElement = document.getElementById("user-list");

  const protocol = window.location.protocol === "https:" ? "wss" : "ws";
  const wsUrl = `${protocol}://${window.location.host}/ws`;
  const socket = new WebSocket(wsUrl);

  socket.onmessage = function (event) {
    const msg = JSON.parse(event.data);
    if (msg.type === "chat") {
      appendMessage(msg.data);
    } else if (msg.type === "userlist") {
      updateUserList(msg.data);
    }
  };

  chatForm.addEventListener("submit", function (e) {
    e.preventDefault();
    if (messageInput.value.trim() === "") return;
    socket.send(messageInput.value);
    messageInput.value = "";
  });

  function appendMessage(message) {
    const div = document.createElement("div");
    div.classList.add("mb-2");

    const time = new Date(message.timestamp).toLocaleTimeString();
    const timeElem = document.createElement("small");
    timeElem.classList.add("text-muted", "me-2");
    timeElem.textContent = time;

    const userElem = document.createElement("strong");
    userElem.textContent = message.user + ":";

    const contentElem = document.createElement("span");
    contentElem.textContent = " " + message.content;

    div.appendChild(timeElem);
    div.appendChild(userElem);
    div.appendChild(contentElem);

    chatBox.appendChild(div);
    chatBox.scrollTop = chatBox.scrollHeight;
  }

  function updateUserList(users) {
    userListElement.innerHTML = "";
    users.forEach((user) => {
      const li = document.createElement("li");
      li.classList.add("list-group-item");
      li.textContent = user;
      userListElement.appendChild(li);
    });
  }
});
