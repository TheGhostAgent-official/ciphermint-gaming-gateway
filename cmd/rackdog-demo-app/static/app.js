// Simple in-browser state to mirror the CipherMint™ model at a high level
const state = {
  player: null,
  balance: 0,
  pricePerTokenUSD: 0.01, // demo valuation: $0.01 per RACKDOG™
  platform: null,
};

// DOM helpers
const $ = (id) => document.getElementById(id);
const stepButtons = document.querySelectorAll(".step-pill");
const screens = {
  welcome: $("screen-welcome"),
  account: $("screen-account"),
  card: $("screen-card"),
  wallet: $("screen-wallet"),
  platform: $("screen-platform"),
  withdraw: $("screen-withdraw"),
};

function setActiveStep(step) {
  stepButtons.forEach((btn) =>
    btn.classList.toggle("active", btn.dataset.step === step)
  );

  Object.entries(screens).forEach(([key, el]) => {
    el.classList.toggle("screen-active", key === step);
  });
}

// Navigate via step pills
stepButtons.forEach((btn) => {
  btn.addEventListener("click", () => {
    const step = btn.dataset.step;
    setActiveStep(step);
  });
});

// Welcome screen
$("btn-start").addEventListener("click", () => {
  setActiveStep("account");
});

// Account creation
$("form-account").addEventListener("submit", (e) => {
  e.preventDefault();
  const gamertag = $("acct-gamertag").value.trim();
  const email = $("acct-email").value.trim();
  const phone = $("acct-phone").value.trim();
  const password = $("acct-password").value.trim();
  const terms = $("acct-terms").checked;

  if (!gamertag || !email || !phone || !password || !terms) {
    alert("Please complete all fields and accept the terms.");
    return;
  }

  state.player = {
    gamertag,
    email,
    phone,
  };

  setActiveStep("card");
});

// Card setup
$("form-card").addEventListener("submit", (e) => {
  e.preventDefault();

  const cardNumber = $("card-number").value.trim();
  const expiry = $("card-expiry").value.trim();
  const cvv = $("card-cvv").value.trim();
  const zip = $("card-zip").value.trim();

  if (!cardNumber || !expiry || !cvv || !zip) {
    alert("Please fill out all card fields.");
    return;
  }

  // Simulate a successful card tokenization + signup bonus
  state.balance = 100; // 100 RACKDOG™ sign-up bonus

  // Update wallet screen
  $("wallet-player").textContent = state.player.gamertag;
  $("wallet-balance").textContent = state.balance.toString();
  $("wallet-value").textContent = formatUSD(
    state.balance * state.pricePerTokenUSD
  );
  $("withdraw-available").value = state.balance.toString();

  setActiveStep("wallet");
});

// Earnings breakdown toggle
$("btn-view-earnings").addEventListener("click", () => {
  const panel = $("earnings-panel");
  panel.hidden = !panel.hidden;
});

// Platform linking
const platformButtons = document.querySelectorAll(".platform-btn");
platformButtons.forEach((btn) => {
  btn.addEventListener("click", () => {
    const platform = btn.dataset.platform;
    state.platform = platform;

    $("platform-status").hidden = false;
    $("platform-text").textContent = `Player ${
      state.player ? state.player.gamertag : "RackDog™ User"
    } is now linked to ${platform}. In production, this connection would route events from ${platform} titles into the CipherMint™ Gaming Gateway for rewards.`;

    setActiveStep("platform");
  });
});

// Withdraw simulation
$("form-withdraw").addEventListener("submit", (e) => {
  e.preventDefault();
  const amount = parseInt($("withdraw-amount").value, 10);

  if (isNaN(amount) || amount <= 0) {
    alert("Enter a valid RACKDOG™ amount.");
    return;
  }

  if (amount > state.balance) {
    alert("You cannot withdraw more than your available RACKDOG™ balance.");
    return;
  }

  state.balance -= amount;
  $("wallet-balance").textContent = state.balance.toString();
  $("wallet-value").textContent = formatUSD(
    state.balance * state.pricePerTokenUSD
  );
  $("withdraw-available").value = state.balance.toString();

  const usd = amount * state.pricePerTokenUSD;
  const msg = `Demo: A withdrawal of ${amount} RACKDOG™ (~${formatUSD(
    usd
  )}) has been queued for payout to the saved card. In a live deployment, this would be processed via CipherMint™ + a compliant payout provider.`;

  $("withdraw-result").hidden = false;
  $("withdraw-message").textContent = msg;

  setActiveStep("withdraw");
});

function formatUSD(value) {
  return `$${value.toFixed(2)}`;
}
