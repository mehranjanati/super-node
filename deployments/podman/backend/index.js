require('dotenv').config();
const express = require('express');
const cors = require('cors');
const { AccessToken } = require('livekit-server-sdk');
const { randomUUID } = require('crypto');

const app = express();
const port = process.env.PORT || 3000;

app.use(cors());
app.use(express.json());

app.post('/token', async (req, res) => {
    try {
        const { roomName, participantName } = req.body;
        
        const room = roomName || 'lobby';
        const participant = participantName || randomUUID();

        const apiKey = process.env.LIVEKIT_API_KEY;
        const apiSecret = process.env.LIVEKIT_API_SECRET;
        const wsUrl = process.env.LIVEKIT_URL;

        if (!apiKey || !apiSecret || !wsUrl) {
            return res.status(500).json({ error: 'Server misconfigured' });
        }

        const at = new AccessToken(apiKey, apiSecret, {
            identity: participant,
        });

        at.addGrant({
            roomJoin: true,
            room: room,
            canPublish: true,
            canSubscribe: true,
        });

        const token = await at.toJwt();

        res.json({
            token,
            serverUrl: wsUrl
        });
    } catch (error) {
        console.error('Error generating token:', error);
        res.status(500).json({ error: 'Failed to generate token' });
    }
});

app.listen(port, () => {
    console.log(`Backend server running on port ${port}`);
});
