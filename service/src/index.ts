import express from "express";
import cors from "cors";

type Coin = {
  denom: string;
  amount: string; // micro units as string
};

type WalletRecord = {
  playerId: string;
  walletAddress: string;
  label: string;
  balances: Coin[];
};

// Demo mapping: player → wallet + balances
const DEMO_WALLETS: WalletRecord[] = [
  {
    playerId: "player-main",
    walletAddress: "ciphermint1demoaddressxyz",
    label: "Main CipherMint Wallet",
    balances: [
      { denom: "ucmint", amount: "125000000" }, // 125 CMINT
      { denom: "urackd", amount: "84500000" },  // 84.5 RACKD
      { denom: "ugame",  amount: "56000000" },  // 56 GAME
    ],
  },
  {
    playerId: "player-gamer",
    walletAddress: "ciphermintgamerxyz000001",
    label: "Gamer Earnings Wallet",
    balances: [
      { denom: "ucmint", amount: "75000000" },
      { denom: "urackd", amount: "12500000" },
      { denom: "ugame",  amount: "98000000" },
    ],
  },
  {
    playerId: "player-creator",
    walletAddress: "ciphermintcreatorxyz0001",
    label: "Creator Royalty Wallet",
    balances: [
      { denom: "ucmint", amount: "30000000" },
      { denom: "urackd", amount: "45000000" },
      { denom: "ugame",  amount: "22000000" },
    ],
  },
];

const app = express();
app.use(cors());
app.use(express.json());

const PORT = Number(process.env.PORT ?? 4000);

// Health check
app.get("/health", (_req, res) => {
  res.json({
    ok: true,
    service: "ciphermint-gaming-gateway",
    port: PORT,
  });
});

// Fetch mapped wallet for a player
app.get("/player/:playerId/wallet", (req, res) => {
  const { playerId } = req.params;
  const record = DEMO_WALLETS.find((w) => w.playerId === playerId);

  if (!record) {
    return res.status(404).json({
      ok: false,
      error: "PLAYER_NOT_FOUND",
      playerId,
    });
  }

  return res.json({
    ok: true,
    playerId: record.playerId,
    walletAddress: record.walletAddress,
    label: record.label,
    balances: record.balances,
  });
});

// Link or update a player → wallet mapping
app.post("/link-wallet", (req, res) => {
  const { playerId, walletAddress, label } = req.body ?? {};

  if (!playerId || !walletAddress) {
    return res.status(400).json({
      ok: false,
      error: "MISSING_FIELDS",
      details: "playerId and walletAddress are required",
    });
  }

  let record = DEMO_WALLETS.find((w) => w.playerId === playerId);

  if (record) {
    // Update existing mapping
    record.walletAddress = walletAddress;
    if (label) record.label = label;
  } else {
    // Create new mapping with zero balances for now
    record = {
      playerId,
      walletAddress,
      label: label || "Linked Wallet",
      balances: [],
    };
    DEMO_WALLETS.push(record);
  }

  return res.json({
    ok: true,
    playerId: record.playerId,
    walletAddress: record.walletAddress,
    label: record.label,
    balances: record.balances,
  });
});

app.listen(PORT, () => {
  console.log(`🎮 Gaming Gateway running on port ${PORT}`);
});
