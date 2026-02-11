function setMessage(el, text, type) {
  if (!el) return;
  el.textContent = text || "";
  el.classList.remove("success", "error");
  if (type) {
    el.classList.add(type);
  }
}

function getToken() {
  return localStorage.getItem("jwt_token") || "";
}

function setToken(token) {
  if (token) {
    localStorage.setItem("jwt_token", token);
  }
}

function clearToken() {
  localStorage.removeItem("jwt_token");
}

async function apiRequest(path, options) {
  const token = getToken();
  const headers = options.headers || {};
  if (token) {
    headers.Authorization = "Bearer " + token;
  }
  headers["Content-Type"] = "application/json";
  const response = await fetch(path, { ...options, headers });
  const data = await response.json().catch(() => ({}));
  return { ok: response.ok, status: response.status, data };
}

function initLogin() {
  const form = document.getElementById("login-form");
  if (!form) return;
  const message = document.getElementById("login-message");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    setMessage(message, "", "");

    const payload = {
      username: form.username.value.trim(),
      password: form.password.value,
    };

    const result = await apiRequest("/auth/login", {
      method: "POST",
      body: JSON.stringify(payload),
    });

    if (result.ok) {
      setToken(result.data.token || "");
      setMessage(message, "Login успешен. Токен сохранен.", "success");
      form.reset();
      return;
    }

    setMessage(message, result.data.error || "Ошибка входа", "error");
  });
}

function initRegister() {
  const form = document.getElementById("register-form");
  if (!form) return;
  const message = document.getElementById("register-message");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    setMessage(message, "", "");

    const payload = {
      name: form.name.value.trim(),
      surname: form.surname.value.trim(),
      email: form.email.value.trim(),
      username: form.username.value.trim(),
      password: form.password.value,
    };

    const result = await apiRequest("/auth/register", {
      method: "POST",
      body: JSON.stringify(payload),
    });

    if (result.ok) {
      setMessage(message, "Регистрация успешна. Теперь войдите.", "success");
      form.reset();
      return;
    }

    setMessage(message, result.data.error || "Ошибка регистрации", "error");
  });
}

function initClubs() {
  const tbody = document.getElementById("clubs-table-body");
  if (!tbody) return;
  const message = document.getElementById("clubs-message");

  (async () => {
    setMessage(message, "", "");
    const result = await apiRequest("/api/clubs", { method: "GET" });

    if (!result.ok) {
      setMessage(message, result.data.error || "Не удалось загрузить клубы", "error");
      return;
    }

    const clubs = Array.isArray(result.data) ? result.data : [];
    if (clubs.length === 0) {
      setMessage(message, "Клубы не найдены", "error");
      return;
    }

    tbody.innerHTML = "";
    clubs.forEach((club) => {
      const row = document.createElement("tr");
      row.innerHTML = `
        <td>${club.id}</td>
        <td>${club.name}</td>
        <td>${club.city}</td>
        <td>${club.address}</td>
      `;
      tbody.appendChild(row);
    });
  })();
}

function initBooking() {
  const form = document.getElementById("booking-form");
  if (!form) return;
  const message = document.getElementById("booking-message");

  form.addEventListener("submit", async (e) => {
    e.preventDefault();
    setMessage(message, "", "");

    const start = form.start_time.value;
    const end = form.end_time.value;
    if (!start || !end) {
      setMessage(message, "Укажите время начала и окончания", "error");
      return;
    }

    const payload = {
      pc_id: Number(form.pc_id.value),
      user_id: Number(form.user_id.value),
      start_time: new Date(start).toISOString(),
      end_time: new Date(end).toISOString(),
    };

    const result = await apiRequest("/api/bookings", {
      method: "POST",
      body: JSON.stringify(payload),
    });

    if (result.ok) {
      setMessage(message, "Бронирование создано", "success");
      form.reset();
      return;
    }

    setMessage(message, result.data.error || "Ошибка бронирования", "error");
  });
}

function initLogout() {
  const link = document.getElementById("logout-link") || document.getElementById("logout-button");
  if (!link) return;

  link.addEventListener("click", async (e) => {
    e.preventDefault();
    await apiRequest("/api/logout", { method: "POST" });
    clearToken();
    window.location = "/login";
  });
}

document.addEventListener("DOMContentLoaded", () => {
  initLogin();
  initRegister();
  initClubs();
  initBooking();
  initLogout();
});
