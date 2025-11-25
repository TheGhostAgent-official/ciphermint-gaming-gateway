import React, { useState } from 'react';
import './App.css';

type Step = 'welcome' | 'create' | 'card' | 'network' | 'wallet';

function App() {
  const [step, setStep] = useState<Step>('welcome');
  const [gamerTag, setGamerTag] = useState('');
  const [network, setNetwork] = useState<'PlayStation' | 'Xbox' | 'Nintendo' | ''>('');
  const [hasCard, setHasCard] = useState(false);
  const [balance, setBalance] = useState(0);

  const handleStart = () => {
    setStep('create');
  };

  const handleCreateWallet = () => {
    if (!gamerTag.trim()) return;
    // Demo: instantly give 100 RACKDOG
    setBalance(100);
    setStep('card');
  };

  const handleAddCard = () => {
    setHasCard(true);
    setStep('network');
  };

  const handleConnectNetwork = () => {
    if (!network) return;
    setStep('wallet');
  };

  const resetFlow = () => {
    setStep('welcome');
    setGamerTag('');
    setNetwork('');
    setHasCard(false);
    setBalance(0);
  };

  const neonBackground: React.CSSProperties = {
    minHeight: '100vh',
    margin: 0,
    padding: 0,
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
    background:
      'radial-gradient(circle at top, rgba(0, 255, 255, 0.35), transparent 60%), ' +
      'radial-gradient(circle at bottom, rgba(255, 0, 255, 0.4), transparent 60%), ' +
      'radial-gradient(circle at left, rgba(255, 140, 0, 0.25), transparent 55%), ' +
      'radial-gradient(circle at right, rgba(0, 140, 255, 0.25), transparent 55%), ' +
      'linear-gradient(135deg, #040516 0%, #05071f 35%, #050312 70%, #080016 100%)',
    color: '#f7f4ff',
    fontFamily: 'system-ui, -apple-system, BlinkMacSystemFont, "SF Pro Text", sans-serif',
  };

  const cardStyle: React.CSSProperties = {
    width: '100%',
    maxWidth: 420,
    borderRadius: 24,
    padding: '26px 24px 22px',
    background:
      'linear-gradient(135deg, rgba(10, 10, 30, 0.98), rgba(7, 9, 40, 0.98))',
    boxShadow:
      '0 18px 50px rgba(0, 0, 0, 0.85), 0 0 40px rgba(0, 255, 255, 0.22), 0 0 70px rgba(255, 0, 255, 0.25)',
    border: '1px solid rgba(0, 255, 255, 0.18)',
    backdropFilter: 'blur(18px)',
    WebkitBackdropFilter: 'blur(18px)',
  };

  const headerRow: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'space-between',
    marginBottom: 18,
    gap: 12,
  };

  const brandRow: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    gap: 12,
  };

  const logoChip: React.CSSProperties = {
    width: 42,
    height: 42,
    borderRadius: 12,
    background:
      'radial-gradient(circle at 0% 0%, #00f0ff 0%, #0066ff 30%, #3b0d78 70%, #0b031a 100%)',
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    boxShadow: '0 0 24px rgba(0, 240, 255, 0.65)',
    overflow: 'hidden',
  };

  const rackdogMini: React.CSSProperties = {
    width: 42,
    height: 42,
    borderRadius: '50%',
    border: '2px solid rgba(255, 170, 0, 0.9)',
    backgroundImage: 'url("/rackdog-coin.png")',
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    boxShadow: '0 0 24px rgba(255, 170, 0, 0.75)',
  };

  const chipText: React.CSSProperties = {
    display: 'flex',
    flexDirection: 'column',
    gap: 2,
  };

  const titleText: React.CSSProperties = {
    fontSize: 22,
    fontWeight: 700,
    letterSpacing: 0.2,
  };

  const subtitleText: React.CSSProperties = {
    fontSize: 11,
    textTransform: 'uppercase',
    letterSpacing: 1.3,
    color: 'rgba(221, 228, 255, 0.72)',
  };

  const badgeRow: React.CSSProperties = {
    display: 'flex',
    alignItems: 'center',
    gap: 8,
  };

  const poweredBadge: React.CSSProperties = {
    fontSize: 10,
    textTransform: 'uppercase',
    letterSpacing: 1.2,
    padding: '4px 10px',
    borderRadius: 999,
    border: '1px solid rgba(0, 255, 255, 0.45)',
    background:
      'linear-gradient(90deg, rgba(0, 255, 255, 0.15), rgba(255, 0, 255, 0.15))',
    color: 'rgba(214, 245, 255, 0.96)',
    whiteSpace: 'nowrap',
  };

  const chipPill: React.CSSProperties = {
    padding: '3px 10px',
    borderRadius: 999,
    fontSize: 10,
    textTransform: 'uppercase',
    letterSpacing: 1.1,
    background: 'rgba(6, 12, 60, 0.9)',
    border: '1px solid rgba(0, 255, 255, 0.25)',
    color: 'rgba(208, 230, 255, 0.95)',
    whiteSpace: 'nowrap',
  };

  const mainContent: React.CSSProperties = {
    marginTop: 10,
  };

  const h1Style: React.CSSProperties = {
    fontSize: 19,
    fontWeight: 700,
    marginBottom: 4,
  };

  const h2Style: React.CSSProperties = {
    fontSize: 13,
    fontWeight: 500,
    color: 'rgba(220, 232, 255, 0.88)',
    marginBottom: 10,
  };

  const copyStyle: React.CSSProperties = {
    fontSize: 12,
    lineHeight: 1.5,
    color: 'rgba(198, 213, 255, 0.9)',
    marginBottom: 18,
  };

  const inputStyle: React.CSSProperties = {
    width: '100%',
    borderRadius: 999,
    border: '1px solid rgba(0, 255, 255, 0.4)',
    padding: '10px 14px',
    fontSize: 13,
    backgroundColor: 'rgba(3, 6, 30, 0.9)',
    color: '#f7f4ff',
    outline: 'none',
    marginBottom: 10,
  };

  const selectStyle: React.CSSProperties = {
    ...inputStyle,
    paddingRight: 28,
    appearance: 'none',
  };

  const primaryButton: React.CSSProperties = {
    width: '100%',
    border: 'none',
    borderRadius: 999,
    padding: '11px 16px',
    marginTop: 6,
    background:
      'linear-gradient(90deg, #00f0ff 0%, #00a3ff 22%, #7a2cff 62%, #ff00d4 100%)',
    color: '#040416',
    fontWeight: 700,
    fontSize: 13,
    letterSpacing: 0.4,
    textTransform: 'uppercase',
    boxShadow:
      '0 12px 32px rgba(0, 0, 0, 0.75), 0 0 40px rgba(0, 240, 255, 0.65)',
    cursor: 'pointer',
  };

  const ghostButton: React.CSSProperties = {
    width: '100%',
    borderRadius: 999,
    padding: '10px 16px',
    background: 'transparent',
    border: '1px solid rgba(0, 255, 255, 0.25)',
    color: 'rgba(210, 225, 255, 0.9)',
    fontSize: 12,
    marginTop: 10,
    cursor: 'pointer',
  };

  const labelStyle: React.CSSProperties = {
    fontSize: 11,
    textTransform: 'uppercase',
    letterSpacing: 1.2,
    color: 'rgba(168, 188, 255, 0.85)',
    marginBottom: 4,
  };

  const pillRow: React.CSSProperties = {
    display: 'flex',
    flexWrap: 'wrap',
    gap: 8,
    marginBottom: 10,
  };

  const pill: (active: boolean) => React.CSSProperties = (active) => ({
    padding: '7px 12px',
    borderRadius: 999,
    fontSize: 11,
    border: active
      ? '1px solid rgba(0, 255, 255, 0.9)'
      : '1px solid rgba(77, 99, 150, 0.9)',
    background: active
      ? 'radial-gradient(circle at 0% 0%, rgba(0, 255, 255, 0.3), rgba(10, 15, 50, 0.95))'
      : 'rgba(6, 10, 40, 0.95)',
    color: 'rgba(214, 228, 255, 0.98)',
    cursor: 'pointer',
  });

  const footerNote: React.CSSProperties = {
    fontSize: 10,
    lineHeight: 1.5,
    color: 'rgba(158, 178, 230, 0.85)',
    marginTop: 18,
  };

  const balanceBadge: React.CSSProperties = {
    fontSize: 11,
    padding: '6px 10px',
    borderRadius: 999,
    background:
      'linear-gradient(120deg, rgba(255, 170, 0, 0.14), rgba(255, 70, 0, 0.16))',
    border: '1px solid rgba(255, 190, 70, 0.9)',
    color: 'rgba(255, 226, 177, 0.98)',
    display: 'inline-flex',
    alignItems: 'center',
    gap: 6,
  };

  const ciphermintTag: React.CSSProperties = {
    fontSize: 10,
    textTransform: 'uppercase',
    letterSpacing: 1.1,
    color: 'rgba(151, 221, 255, 0.95)',
  };

  const renderStepBody = () => {
    switch (step) {
      case 'welcome':
        return (
          <>
            <div style={h1Style}>RackDawg™ Gaming Wallet</div>
            <div style={h2Style}>Turn playtime into real value.</div>
            <p style={copyStyle}>
              Create a CipherMint-powered wallet, collect a{' '}
              <strong>100&nbsp;RACKDOG™</strong> sign-up bonus, and wire your
              value directly into the games and networks you already play on.
            </p>
            <button style={primaryButton} onClick={handleStart}>
              Start Demo • Create Wallet
            </button>
            <p style={footerNote}>
              Demo only — in a live integration, balances are backed by the
              CipherMint™ chain and can be spent in-game or withdrawn to a
              player&apos;s card.
            </p>
          </>
        );
      case 'create':
        return (
          <>
            <div style={h1Style}>Create your RackDawg™ wallet</div>
            <div style={h2Style}>
              Pick a gamer tag to mint your demo wallet and sign-up bonus.
            </div>
            <label style={labelStyle}>Gamer tag</label>
            <input
              style={inputStyle}
              placeholder="GhostPlayer, RackSniper, etc."
              value={gamerTag}
              onChange={(e) => setGamerTag(e.target.value)}
            />
            <button style={primaryButton} onClick={handleCreateWallet}>
              Mint Wallet • Get 100 RACKDOG
            </button>
            <button style={ghostButton} onClick={resetFlow}>
              ← Back to welcome
            </button>
          </>
        );
      case 'card':
        return (
          <>
            <div style={h1Style}>Add a demo card</div>
            <div style={h2Style}>
              In production, this connects to a real debit or credit card.
              Here, it&apos;s a safe demo.
            </div>
            <label style={labelStyle}>Demo card number</label>
            <input
              style={inputStyle}
              placeholder="4242 4242 4242 4242"
              value="4242 4242 4242 4242"
              readOnly
            />
            <label style={labelStyle}>Expiry • CVC</label>
            <input
              style={inputStyle}
              placeholder="12 / 28 • 123"
              value="12 / 28 • 123"
              readOnly
            />
            <button style={primaryButton} onClick={handleAddCard}>
              Continue • Link demo card
            </button>
            <button style={ghostButton} onClick={resetFlow}>
              ← Start over
            </button>
          </>
        );
      case 'network':
        return (
          <>
            <div style={h1Style}>Connect your gaming network</div>
            <div style={h2Style}>
              Choose where this wallet should unlock value first.
            </div>
            <label style={labelStyle}>Gaming network</label>
            <div style={pillRow}>
              {(['PlayStation', 'Xbox', 'Nintendo'] as const).map((n) => (
                <button
                  key={n}
                  type="button"
                  style={pill(network === n)}
                  onClick={() => setNetwork(n)}
                >
                  {n}
                </button>
              ))}
            </div>
            <button style={primaryButton} onClick={handleConnectNetwork}>
              Connect &amp; Finish
            </button>
            <button style={ghostButton} onClick={resetFlow}>
              ← Start over
            </button>
          </>
        );
      case 'wallet':
        return (
          <>
            <div style={h1Style}>You&apos;re wired in.</div>
            <div style={h2Style}>
              Your RackDawg™ demo wallet is live on CipherMint™.
            </div>
            <p style={copyStyle}>
              This screen represents the state we expose to games, consoles, and
              studios: balances, linked cards, and networks are all available by
              API through the CipherMint™ Gaming Gateway.
            </p>
            <div style={{ marginBottom: 12 }}>
              <span style={balanceBadge}>
                <span
                  style={{
                    width: 8,
                    height: 8,
                    borderRadius: '50%',
                    background:
                      'radial-gradient(circle, #ffcc4d 0%, #ff8c00 60%, #ff0050 100%)',
                    display: 'inline-block',
                  }}
                />
                <span>
                  {balance} RACKDOG • {gamerTag || 'GhostPlayer'}
                </span>
              </span>
            </div>
            <div style={{ fontSize: 11, marginBottom: 8 }}>
              Linked network:{' '}
              <strong>{network || 'Not set in this demo'}</strong>
            </div>
            <div style={{ fontSize: 11, marginBottom: 16 }}>
              Demo card: <strong>{hasCard ? 'Connected' : 'Not connected'}</strong>
            </div>
            <button style={ghostButton} onClick={resetFlow}>
              Run the flow again
            </button>
            <p style={footerNote}>
              In production, this wallet can receive on-chain rewards, route
              payouts to cards, and sync balances into any supported game
              engine, marketplace, or console network.
            </p>
          </>
        );
      default:
        return null;
    }
  };

  return (
    <div style={neonBackground}>
      <div style={cardStyle}>
        <div style={headerRow}>
          <div style={brandRow}>
            <div style={logoChip}>
              {/* CipherMint neon hex logo */}
              <div
                style={{
                  width: 26,
                  height: 26,
                  borderRadius: 8,
                  border: '2px solid rgba(0, 255, 255, 0.9)',
                  boxShadow: '0 0 20px rgba(0, 255, 255, 0.9)',
                  backgroundImage: 'url("/ciphermint-logo.png")',
                  backgroundSize: 'cover',
                  backgroundPosition: 'center',
                }}
              />
            </div>
            <div style={chipText}>
              <span style={titleText}>RackDawg™ Gaming Wallet</span>
              <span style={subtitleText}>CipherMint™ Studios • Demo</span>
            </div>
          </div>
          <div style={rackdogMini} />
        </div>

        <div style={badgeRow}>
          <span style={poweredBadge}>Powered by CipherMint™ Gaming Gateway</span>
          <span style={chipPill}>RackDawg™ • RACKDOG</span>
        </div>

        <div style={mainContent}>{renderStepBody()}</div>

        <div style={{ marginTop: 14, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          <span style={ciphermintTag}>CipherMint • Utility-first token rails</span>
        </div>
      </div>
    </div>
  );
}

export default App;
