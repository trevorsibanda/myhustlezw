import store from "store2"


var usdFormat = Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'USD',
        minimumFractionDigits: 2
})

var zwlFormat = Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: 'ZWL',
        minimumFractionDigits: 2
})
    
function formatUSD(amount) {
    return usdFormat.format(amount)
}

function formatZWL(amount) {
    return zwlFormat.format(amount)
}


function formatUSDToZWL(amount_usd) {
    var usdToZWLRate = store.get('currency_rate', 700)
    return formatZWL(amount_usd * usdToZWLRate)
}

function formatMoney(curreny, amount) {
    if(curreny === 'USD') {
        return formatUSD(amount)
    }
    return formatZWL(amount)
}

function USDPrice(props) {
    return <span>{formatUSD(props.amount)}</span>
}

function ZWLPrice(props) {
    return <span>{props.usd ? formatUSDToZWL(props.usd) : formatZWL(props.amount)}</span>
}

let components = {
    USD: USDPrice,
    ZWL: ZWLPrice,
    formatUSDToZWL: formatUSDToZWL,
    formatUSD: formatUSD,
    formatZWL: formatZWL,
    format:  formatMoney,
}

export default components;