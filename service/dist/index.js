"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const cors_1 = __importDefault(require("cors"));
const app = (0, express_1.default)();
const PORT = process.env.PORT || 4000;
// Basic middleware
app.use((0, cors_1.default)());
app.use(express_1.default.json());
const DEMO_PLAYER_ID = "demo-player-1";
// ---- Demo data (CipherMint × Apple vibe) ----
const demoCoins = [
    {
        symbol: "CIPH",
        label: "CipherMint",
        balance: 128430.12,
        usdValue: 1284.3,
    },
    {
        symbol: "RACKD",
        label: "RackDawg",
        balance: 52900,
        usdValue: 529.0,
    },
    {
        symbol: "APPLΞ",
        label: "Apple Energy",
        balance: 10000,
        usdValue: 300.0,
    },
];
const demoTransactions = [
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
const demoSnapshot = {
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
        appleGameCenterId: appleGameCenterId !== null && appleGameCenterId !== void 0 ? appleGameCenterId : null,
    });
});
// ---- Start server ----
app.listen(PORT, () => {
    // eslint-disable-next-line no-console
    console.log(`CipherMint Gaming Gateway running on port ${PORT} (demo player endpoints enabled)`);
});
