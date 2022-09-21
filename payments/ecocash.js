const { Paynow } = require("paynow");
const LRU = require("./lrucache")
const fetch = require("node-fetch")

function initiateTransaction(mobile_no, amount, ref, title, email){
    let paynow = new Paynow(process.env.PAYNOW_INTEGRATION_ID, process.env.PAYNOW_INTEGRATION_KEY);

    console.log(process.env.PAYNOW_INTEGRATION_ID, process.env.PAYNOW_INTEGRATION_KEY);

    paynow.resultUrl = process.env.PAYMENT_API_CALLBACK_ENDPOINT + "/ecocash" 
    

    let payment = paynow.createPayment(ref, email);
    payment.add(title, amount, 1)
    //console.log(title, amount);
    //console.log(payment);
    //payment.authEmail = 'sibandatrevor@gmail.com'
    return paynow.sendMobile(payment, '0'+mobile_no, 'ecocash')
}

function pollTransaction(payment_id) {
    let response = LRU.get(payment_id)
    //do something
    return fetch(response.pollUrl).then(res => res.text())
}

exports.poll = pollTransaction

exports.init = initiateTransaction