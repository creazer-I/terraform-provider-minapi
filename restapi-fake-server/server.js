const express = require('express');
const fs = require('fs');
const app = express();
app.use(express.json());

// Middleware for checking the token
app.use((req, res, next) => {
    const token = req.header('Authorization');
    if (!token) {
        return res.status(401).send('Access denied. No token provided.');
    }
    if (token !== 'fakeservertoken') {
        return res.status(401).send('Access denied. Invalid token.');
    }
    next();
});

app.post('/test', (req, res) => {
    // Append payload data to 'payload.json'
    fs.appendFile('payload.json', JSON.stringify(req.body, null, 2) + ','+'\n', (err) => {
        if (err) {
            console.error(err);
            res.status(500).send('An error occurred while writing to file');
        } else {
            res.send('Payload appended to file successfully');
        }
    });
});

app.listen(3000, () => console.log('Server is running on port 3000'));
