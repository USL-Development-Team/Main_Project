package templates

const LoginFormHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>%s</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
        }
        .login-container {
            background-color: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
            max-width: 400px;
            width: 100%%;
            text-align: center;
        }
        h2 {
            color: #333;
            margin-bottom: 20px;
        }
        .usl-header {
            color: #1a365d;
            margin-bottom: 10px;
            font-size: 24px;
            font-weight: bold;
        }
        .discord-btn {
            display: block;
            width: 100%%;
            padding: 15px;
            background-color: #5865F2;
            color: white;
            text-decoration: none;
            border-radius: 5px;
            font-size: 16px;
            font-weight: bold;
            transition: background-color 0.3s;
            text-align: center;
            box-sizing: border-box;
        }
        .discord-btn:hover {
            background-color: #4752C4;
        }
        .info {
            color: #666;
            margin-bottom: 30px;
            font-size: 14px;
        }
        .error {
            color: red;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="login-container">
        <h2 class="%s">%s</h2>
        <div class="info">%s</div>
        %s
        <a href="%s" class="discord-btn">
            🎮 Sign in with Discord
        </a>
    </div>
</body>
</html>`

const AuthCallbackHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Processing Login...</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            display: flex;
            justify-content: center;
            align-items: center;
            min-height: 100vh;
            background-color: #f5f5f5;
            margin: 0;
        }
        .processing {
            text-align: center;
            background: white;
            padding: 40px;
            border-radius: 8px;
            box-shadow: 0 2px 10px rgba(0,0,0,0.1);
        }
        .spinner {
            border: 4px solid #f3f3f3;
            border-top: 4px solid #3498db;
            border-radius: 50%%;
            width: 30px;
            height: 30px;
            animation: spin 1s linear infinite;
            margin: 20px auto;
        }
        @keyframes spin {
            0%% { transform: rotate(0deg); }
            100%% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="processing">
        <h2>Processing Login...</h2>
        <div class="spinner"></div>
        <p>Please wait while we complete your authentication.</p>
    </div>
    <script>
        // Extract tokens from URL fragment (Supabase returns them in hash)
        const hash = window.location.hash.substring(1);
        const params = new URLSearchParams(hash);
        const accessToken = params.get('access_token');
        const refreshToken = params.get('refresh_token');
        
        if (accessToken) {
            // Send tokens to server for validation and session setup
            fetch('/auth/process', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    access_token: accessToken,
                    refresh_token: refreshToken 
                })
            }).then(response => {
                if (response.ok) {
                    // Successful authentication - redirect to final destination
                    window.location.href = '%s';
                } else {
                    // Authentication failed
                    window.location.href = '/login?error=unauthorized';
                }
            }).catch(error => {
                console.error('Auth error:', error);
                window.location.href = '/login?error=invalid';
            });
        } else {
            // No token found - redirect to login with error
            window.location.href = '/login?error=invalid';
        }
    </script>
</body>
</html>`
