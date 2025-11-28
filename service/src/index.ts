import express from "express";
import cors from "cors";

const app = express();
const PORT = process.env.PORT || 4000;

// Basic middleware
app.use(cors());
app.use(express.json());

// ---- Demo data types ----
type Coin = {
  symbol: string;
  label: string;
  balance: number;
  usdValue: number;
};

type Transaction = {
  id: string;
  type: "earn" | "spend" | "airdrop";
  amount: number;
  coin: string;
  timestamp: string;
  description: string;
};

type PlayerWalletSnapshot = {
  playerId: string;
  address: string;
  totalUsdValue: number;
  coins: Coin[];
  recentTransactions: Transaction[];
  lifetimeEarningsUsd: number;
};

const DEMO_PLAYER_ID = "demo-player-1";

// ---- Demo data (CipherMint × Apple vibe) ----
const demoCoins: Coin[] = [
  {
    symbol: "CIPH",
    label: "CipherMint",
    balance: 128_430.12,
    usdValue: 1284.3,
  },
  {
    symbol: "RACKD",
    label: "RackDawg",
    balance: 52_900.0,
    usdValue: 529.0,
  },
  {
    symbol: "APPLΞ",
    label: "Apple Energy",
    balance: 10_000,
    usdValue: 300.0,
  },
];

const demoTransactions: Transaction[] = [
  {
    id: "tx_001",
    type: "earn",
    amount: 75,
    coin: "CIPH",
    timestamp: new Date(Date.now() - 1000 * 60 * 10).toISOString(),
    description: "Ranked win bonus",
  },
  {
    id: "tx_002",
    type: "airdrop",
    amount: 250,
    coin: "RACKD",
    timestamp: new Date(Date.now() - 1000 * 60 * 60).toISOString(),
    description: "Creator airdrop",
  },
  {
    id: "tx_003",
    type: "spend",
    amount: -35,
    coin: "CIPH",
    timestamp: new Date(Date.now() - 1000 * 60 * 60 * 5).toISOString(),
    description: "In-game skin purchase",
  },
];

const totalUsd = demoCoins.reduce((sum, c) => sum + c.usdValue, 0);
const lifetimeEarningsUsd = 2100.5;

const demoSnapshot: PlayerWalletSnapshot = {
  playerId: DEMO_PLAYER_ID,
  address: "0xC1PHERM1NT-APPLE-DEMO",
  totalUsdValue: totalUsd,
  coins: demoCoins,
  recentTransactions: demoTransactions,
  lifetimeEarningsUsd,
};

// ---- Root + health ----

// Nice message if someone hits the base URL in a browser
app.get("/", (_req, res) => {
  res.send("CipherMint Gaming Gateway • online");
});

app.get("/health", (_req, res) => {
  res.json({
    status: "ok",
    service: "ciphermint-gaming-gateway",
    playerId: DEMO_PLAYER_ID,
    time: new Date().toISOString(),
  });
});

// ---- Wallet endpoints used by the dashboard ----

app.get("/wallet/balances", (_req, res) => {
  res.json({
    playerId: DEMO_PLAYER_ID,
    coins: demoCoins,
    totalUsdValue: totalUsd,
  });
});

app.get("/wallet/transactions", (_req, res) => {
  res.json({
    playerId: DEMO_PLAYER_ID,
    transactions: demoTransactions,
  });
});

app.get("/player/wallet", (_req, res) => {
  res.json(demoSnapshot);
});

app.post("/player/register", (req, res) => {
  const { appleGameCenterId } = req.body || {};
  res.json({
    playerId: DEMO_PLAYER_ID,
    appleGameCenterId: appleGameCenterId ?? null,
  });
});

// ---- Start server ----
app.listen(PORT, () => {
  // eslint-disable-next-line no-console
  console.log(
    `CipherMint Gaming Gateway running on port ${PORT} (demo player endpoints enabled)`
  );
});
