package gomail

const verificationCodeHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Verification Code</title>
</head>
<body>
    <p>This is your verification code: %s</p>
</body>
</html>
`

const invoiceRentHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Rent Invoice</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }

        .invoice-container {
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #fff;
            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
        }

        h1 {
            font-size: 24px;
            color: #333;
        }

        p {
            font-size: 16px;
            color: #666;
        }

        .total {
            font-size: 18px;
            font-weight: bold;
            color: #333;
        }

        .payment-button {
            display: inline-block;
            background-color: #007BFF;
            color: #fff;
            padding: 10px 20px;
            font-size: 16px;
            border: none;
            cursor: pointer;
            margin-top: 20px;
            text-decoration: none;
        }

        .payment-button:hover {
            background-color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="invoice-container">
        <h1>Rent Invoice</h1>
        <p><strong>Book Title:</strong> %s</p>
        <p><strong>Days of Rent:</strong> %d days</p>
        <p><strong>Amount:</strong> Rp.%2.f</p>
		<p><strong>Token:</strong> %s</p>
		<p><strong>Payment Link:</strong> %s</p>
       
    </div>
</body>
</html>
`
const paymentSuccessHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Payment Success</title>
</head>
<body>
    <h1>Payment Successful</h1>
	<p><strong>Order Id:</strong> %s</p>
	<p><strong>Payment Type:</strong> %s</p>
	<p><strong>Paid At:</strong> %s</p>
    <p>Thank you for your payment. Your transaction was successful.</p>
    <p>Waiting for your admin to response your order.</p>
</body>
</html>
`
