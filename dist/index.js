"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
const express_1 = __importDefault(require("express"));
const cors_1 = __importDefault(require("cors"));
/**
 * Helpers
 */
function addAmounts(a, b) {
    return (BigInt(a) + BigInt(b)).toString();
}
function nowIso() {
    return new Date().toISOString();
}
function makeDemoTxHash(prefix = "demo") {
    return `${prefix}_${Date.now().toString(36)}_${Math.random()
        .toString(36)
        .slice(2, 8)}`;
}
/**
 * Demo chain + wallets
 */
const DEMO_CHAIN_STATUS = {
    chainId: "ciphermint-demo-1",
    latestHeight: 18235,
    latestBlockTime: "2025-11-26T11:14:10Z",
    nodeVersion: "CipherMint demo-node v0.1.0",
};
const DEMO_WALLETS = [
    {
        id: "primary",
        label: "Main CipherMint Wallet",
        address: "ciphermint1demoaddressxyz",
        description: "Primary wallet for ecosystem activity.",
    },
    {
        id: "gamer",
        label: "Gamer Earnings Wallet",
        address: "ciphermint1gamerxyz000001",
        description: "In-game rewards and drops.",
    },
    {
        id: "creator",
        label: "Creator Royalty Wallet",
        address: "ciphermint1creatorxyz00001",
        description: "Payouts from creator economies.",
    },
];
const PRIMARY_ADDRESS = DEMO_WALLETS[0].address;
/**
 * Demo balances & transactions per wallet address
 */
const DEMO_BALANCES_BY_ADDRESS = {
    [PRIMARY_ADDRESS]: [
        { denom: "ucmint", amount: "125000000" },
        { denom: "urackd", amount: "84500000" },
    ],
    [DEMO_WALLETS[1].address]: [
        { denom: "ucmint", amount: "56000000" },
        { denom: "urackd", amount: "25000000" },
    ],
    [DEMO_WALLETS[2].address]: [
        { denom: "ucmint", amount: "30000000" },
    ],
};
const DEMO_TXS_BY_ADDRESS = {
    [PRIMARY_ADDRESS]: [],
    [DEMO_WALLETS[1].address]: [],
    [DEMO_WALLETS[2].address]: [],
};
/**
 * Player → wallet mappings
 */
const PLAYER_WALLETS = new Map();
/**
 * Create / attach a wallet for a player
 */
function createWalletForPlayer(playerId) {
    const sanitized = playerId.replace(/[^a-zA-Z0-9]/g, "").toLowerCase() || "player";
    const walletAddress = `ciphermint1${sanitized.slice(0, 20)}demo`;
    if (!DEMO_BALANCES_BY_ADDRESS[walletAddress]) {
        DEMO_BALANCES_BY_ADDRESS[walletAddress] = [
            { denom: "ucmint", amount: "0" },
            { denom: "urackd", amount: "0" },
        ];
    }
    if (!DEMO_TXS_BY_ADDRESS[walletAddress]) {
        DEMO_TXS_BY_ADDRESS[walletAddress] = [];
    }
    const info = {
        playerId,
        walletAddress,
        label: "Gamer Earnings Wallet",
    };
    PLAYER_WALLETS.set(playerId, info);
    return info;
}
/**
 * Express app
 */
const app = (0, express_1.default)();
app.use((0, cors_1.default)());
app.use(express_1.default.json());
const PORT = Number(process.env.PORT) || 4000;
/**
 * Health
 */
app.get("/health", (_req, res) => {
    res.json({ ok: true, service: "ciphermint-gaming-gateway", time: nowIso() });
});
/**
 * Chain + wallet endpoints (for dashboard)
 */
app.get("/chain/status", (_req, res) => {
    res.json(DEMO_CHAIN_STATUS);
});
app.get("/wallet/:address/balances", (req, res) => {
    const { address } = req.params;
    const balances = DEMO_BALANCES_BY_ADDRESS[address] || [];
    res.json({ address, balances });
});
app.get("/wallet/:address/transactions", (req, res) => {
    const { address } = req.params;
    const txs = DEMO_TXS_BY_ADDRESS[address] || [];
    res.json({ address, transactions: txs });
});
/**
 * Player endpoints
 */
app.post("/players/register", (req, res) => {
    const { playerId } = req.body;
    if (!playerId || typeof playerId !== "string" || !playerId.trim()) {
        return res.status(400).json({ error: "playerId is required" });
    }
    const normalized = playerId.trim();
    let info = PLAYER_WALLETS.get(normalized);
    if (!info) {
        info = createWalletForPlayer(normalized);
    }
    return res.json(info);
});
app.get("/players/:playerId/wallet", (req, res) => {
    const { playerId } = req.params;
    const info = PLAYER_WALLETS.get(playerId);
    if (!info) {
        return res.status(404).json({ error: "Player not registered" });
    }
    const balances = DEMO_BALANCES_BY_ADDRESS[info.walletAddress] || [];
    const txs = DEMO_TXS_BY_ADDRESS[info.walletAddress] || [];
    return res.json({
        playerId: info.playerId,
        walletAddress: info.walletAddress,
        label: info.label,
        balances,
        transactions: txs,
    });
});
app.post("/players/:playerId/earn", (req, res) => {
    const { playerId } = req.params;
    const { denom, amount, source, metadata } = req.body;
    if (!denom || typeof denom !== "string") {
        return res.status(400).json({ error: "denom is required" });
    }
    if (!amount || typeof amount !== "string" || !/^[0-9]+$/.test(amount)) {
        return res
            .status(400)
            .json({ error: "amount must be a numeric string in micro units" });
    }
    const info = PLAYER_WALLETS.get(playerId);
    if (!info) {
        return res.status(404).json({ error: "Player not registered" });
    }
    const walletAddress = info.walletAddress;
    const existing = DEMO_BALANCES_BY_ADDRESS[walletAddress] || [];
    let found = false;
    const updatedBalances = existing.map((coin) => {
        if (coin.denom === denom) {
            found = true;
            return {
                denom,
                amount: addAmounts(coin.amount, amount),
            };
        }
        return coin;
    });
    if (!found) {
        updatedBalances.push({ denom, amount });
    }
    DEMO_BALANCES_BY_ADDRESS[walletAddress] = updatedBalances;
    const tx = {
        hash: makeDemoTxHash("earn"),
        height: "0",
        timestamp: nowIso(),
        status: "Success",
        success: true,
        amount,
        denom,
        to: walletAddress,
        source: source || "game_reward",
        metadata: metadata || {},
    };
    if (!DEMO_TXS_BY_ADDRESS[walletAddress]) {
        DEMO_TXS_BY_ADDRESS[walletAddress] = [];
    }
    DEMO_TXS_BY_ADDRESS[walletAddress].unshift(tx);
    return res.json({
        success: true,
        playerId,
        walletAddress,
        balances: DEMO_BALANCES_BY_ADDRESS[walletAddress],
        transactions: DEMO_TXS_BY_ADDRESS[walletAddress],
    });
});
/**
 * Start server
 */
app.listen(PORT, () => {
    console.log(`🎮 CipherMint Gaming Gateway running on port ${PORT} (demo player endpoints enabled)`);
});
