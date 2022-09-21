import store from "store2"
import ReactGA from "react-ga"

let endpoint = ''



function apiEndpoint(uri) {
    return endpoint + "/api/v1/private"+uri
}

function publicApiEndpoint(uri) {
    return endpoint + "/api/v1/public"+uri
}

function apiCallHandler(resp, endpoint) {
    return resp.json().then(response => {
        if (response.error) {
            if(response.error === "Not authenticated") {
                return Promise.reject({type: 'auth_error', response})
            }
            alert({
                title: "Error",
                text: response.error,
                icon: "warning",
            })
            return Promise.reject({ type: 'api_error', response })
        } else {
            //if(endpoint)
            //    store(endpoint, response)
            return Promise.resolve(response)
        }
    })
}

function apicallErrorHandler(err) {
    if (err.type && err.type === 'api_error') {
        console.log('api_error', err)
        return Promise.reject(err.response)
    } else if (err.type && err.type === 'auth_error') {
        window.logout()
        return Promise.reject(err.response)
    }
    alert({
        title: "Network failure",
        text: "Failed to process request with error \n\n" + err,
    })
    return Promise.reject(err)
}

function apiFetchJson(endpoint, cached) {
    if(cached){
        let resp = store.get(endpoint, null)
        if(resp === null) {
            return Promise.reject({error: "cache_miss"})
        }
        return resp
    }
    return fetch(apiEndpoint(endpoint)).then(resp =>apiCallHandler(resp, endpoint)).catch(apicallErrorHandler)
}

function publicApiFetchJson(endpoint) {
    return fetch(publicApiEndpoint(endpoint)).then(resp =>apiCallHandler(resp, endpoint)).catch(apicallErrorHandler)
}

function profPicURL(user) {
    return user.profile.profile_url
}

function imageURL(image_id, width, height) {
    return '/api/v1/private/image/'+width+'/' +height+ '/' + image_id + '.png'
}


function pollPayment(id, signature, ts) {
    return publicApiFetchJson('/transaction/xhr_poll/' + id + '?ts=' + ts + '&signature=' + signature)
}

function dispatchPayment(purpose, payment, nonce) {
    return fetch(publicApiEndpoint('/transaction/initiate/'+purpose + '?nonce='+nonce), {
        method: 'POST',
        body: JSON.stringify(payment),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

function isUserLoggedIn() {
    let endpoint = "/user"
    return fetch(apiEndpoint(endpoint)).then(resp => resp.json()).then(user => {
        if (user.error && user.error === "Not authenticated"){
            return Promise.resolve(false)
        } else {
            return Promise.resolve(user.logged_in)
        }
    }).catch(_ => {
        return Promise.resolve(false)
    })
}

async function setAvatar(image_id) {
    return apiFetchJson('/user/profile/set_avatar/'+image_id)
}

async function serviceTemplates(cached=false) {
    return apiFetchJson('/config/service_templates', cached)
}

function uploadMediaEndpoint(tpe, role){
    return '/api/v1/private/storage/upload/'+tpe + '/'+role
}

async function getAllCampaigns(filters, cached=false) {
    return apiFetchJson("/campaigns/all", cached)
}

async function getLatestSummary(cached=false) {
    return apiFetchJson("/summary", cached)
}

async function getCurrentUser(cached=false){
    return apiFetchJson("/user", cached)
}

async function getRecentSupportersList(max_recent, skip = 0, cached=false){
    return apiFetchJson("/supporters/recent/"+max_recent+"/"+skip, cached)
}

async function getSupporter(id, cached = true) {
    return apiFetchJson("/supporters/get/"+ id, cached)
}

async function hideSupportActivity(id, access) {
    return apiFetchJson("/supporters/hide/"+ access + "/" + id, false)
}

async function modifyServiceOrderStatus(purpose, id, body) {
    return fetch(apiEndpoint('/supporters/service/'+purpose + '/'+id), {
        method: 'POST',
        body: JSON.stringify(body),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function getRecentCampaignSupportersList(campaign_id, max_recent, skip = 0, cached=false) {
    return apiFetchJson("/campaign/supporters/" + campaign_id +'/'+ max_recent + "/" + skip, cached)
}

async function resendSMSCode(cached=false){
    return apiFetchJson("/security/phone/resend_sms", cached)
}

async function resendEmailCode(cached=false) {
    return apiFetchJson("/security/email/resend_email", cached)
}

async function listAllSubscriptions(filter, cached = false) {
    return apiFetchJson("/subscriptions/"+filter, cached)
}

async function getSubscription(id, cached= false) {
    return apiFetchJson("/subscription/"+ id, cached)
}

async function updatePhone(new_phone) {
    return apiFetchJson("/security/phone/update/"+new_phone, false)
}

async function updateEmail(new_email) {
    return apiFetchJson("/security/email/update/" + new_email, false)
}

async function updateCampaign(campaign) {
    //clear cache
    store.clearAll()
    return fetch(apiEndpoint('/campaign/update/'+campaign._id), {
        method: 'POST',
        body: JSON.stringify(campaign),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function unlockCampaign(unlockCode) {
    store.clearAll()
    return publicApiFetchJson(unlockCode, false)
}

async function updatePassword(old, newp) {
    store.clearAll()
    return fetch(apiEndpoint('/security/password/update' ), {
        method: 'POST',
        body: JSON.stringify({old, 'new': newp}),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function verifyPhone(code) {
    store.clearAll()
    return apiFetchJson("/security/phone/verify/"+ code, false)
}

async function verifyEmail(code) {
    store.clearAll()
    return apiFetchJson("/security/email/verify/" + code, false)
}

async function verifyByPayment(phone) {
    store.clearAll()
    return apiFetchJson("/security/identity/verify_by_payment/"+ phone, false)
}

async function loginUser(params) {
    store.clearAll()
    return fetch(apiEndpoint('/security/login'), {
        method: 'POST',
        body: JSON.stringify(params),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function signupUser(params, target) {
    store.clearAll()
    return fetch(apiEndpoint('/security/signup?target='+target), {
        method: 'POST',
        body: JSON.stringify(params),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function resetPasswordRequest(form) {
    store.clearAll()
    return fetch(apiEndpoint('/security/request_reset_password'), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function processPasswordReset(form) {
    store.clearAll()
    return fetch(apiEndpoint('/security/process_reset_password'), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function updateBasicPageDetails(user) {
    store.clearAll()
    return fetch(apiEndpoint('/user/page/basic'), {
        method: 'POST',
        body: JSON.stringify(user),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function saveBasics(form) {
    store.clearAll()
    return fetch(apiEndpoint('/user/update_basics'), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function userLogout() {
    store.clearAll()
    return apiFetchJson("/auth/logout", false)
}

async function updateNotifications(form) {
    store.clearAll()
    return fetch(apiEndpoint('/user/update_notifications'), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function updatePageConfigurables(form) {
    store.clearAll()
    return fetch(apiEndpoint('/user/update_page_configurables'), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function publishPage(bool) {
    store.clearAll()
    return apiFetchJson("/user/page/publish/"+bool, false)
}

async function getCampaign(id, cached=false) {
    return apiFetchJson("/campaign/detailed/" + id, cached)
}

async function createService(form) {
    store.clearAll()
    return fetch(apiEndpoint('/campaign/new/service'), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function createCampaign(form, tpe) {
    store.clearAll()
    return fetch(apiEndpoint('/campaign/new/'+tpe), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

async function deleteCampaign(id) {
    store.clearAll()
    return apiFetchJson("/campaign/delete/"+id, false)
}


async function getWalletSummary(cached=false){
    return apiFetchJson("/wallet/summary", cached)
}

async function deleteFiles(file) {
    return apiFetchJson("/files/delete/"+ file._id, false)
}

async function getRecentWalletOperations(maxPageSize, cached=false) {
    return apiFetchJson("/wallet/operations/"+maxPageSize, cached, cached)
}

async function payoutDetails(cached=false) {
    return apiFetchJson("/wallet/payout_details", cached)
}

async function getPayment(id, cached=false) {
    return apiFetchJson("/wallet/payment/"+id, cached)
}

async function getRecentPayments(maxPageSize, cached=false) {
    return apiFetchJson("/wallet/my_payments/"+maxPageSize, cached)
}

async function sendWithdrawalRequest(currency) {
    store.clearAll()
    return fetch(apiEndpoint('/wallet/withdrawal_request/'+currency), {
        method: 'POST',
        body: JSON.stringify({nonce: Math.random()}),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

function getCreator(username, cached= true){
   return publicApiFetchJson("/_user/" + username, cached)
}
function publicGetCampaign(username, campaign_id, cached= true){
   return publicApiFetchJson("/_campaign/" + username + "/" + campaign_id, cached)
}

function loadFilteredFeed(view, filter, skip) {
    return apiFetchJson("/feed/"+ view+ "/"+filter + "?skip="+skip, false)
}

function publicGetCampaigns(username, page, cached= false){
   return publicApiFetchJson("/_campaigns/" + username + "/"+page, cached)
}
function getMetrics(username, id, cached= false){
   return publicApiFetchJson("/_campaign/metrics/" + username + "/" + id, cached)
}
function pushMetrics(username, id, metrics) {
    return fetch(apiEndpoint("/_campaign/metrics/" + username + "/" + id ), {
        method: 'POST',
        body: JSON.stringify(metrics),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

function requestFileDownloadLink(username, id){
   return publicApiFetchJson("/_file/download/" + username + "/"+id, false)
}

function requestService(username, campaign_id, form, cached= true){
    return fetch(apiEndpoint("/_service/request/" + username + "/"+ campaign_id), {
        method: 'POST',
        body: JSON.stringify(form),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

function getServiceRequest(request_id){
   return publicApiFetchJson("/_service/request/" + request_id, false)
}

function getSiteConfig(cached= true){
   return publicApiFetchJson("/_site/config", cached)
}

function syncCurrenyRate() {
    return publicApiFetchJson("/_site/currency_sync", false).then(currency => {
        if (!currency.error && currency.usd_to_zwl) {
            store.set('currency_rate', currency.usd_to_zwl)
        }
    })
}

async function updateCashout(currency, details) {
    store.clearAll()
    return fetch(apiEndpoint('/wallet/payout_details/' + currency), {
        method: 'POST',
        body: JSON.stringify(details),
    }).then(apiCallHandler).catch(apicallErrorHandler)
}

const validateEmail = (email) => {
  return String(email)
    .toLowerCase()
    .match(
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
    );
}

const validatePhone = (phone) => {
    return String(phone)
        .toLowerCase()
        .match(
            /^\+?[0-9]{6,}$/
        );
}

function generateDefaultCampaigns() {
    return []
}

let defaults = {
    creator: function (username) {
        return {
                username: username,
                fullname: 'Loading...',
                profile: {
                    description: '...',
                    personal_website: {

                    }
                }
            }
    },
    featured: {
        type: 'image',
        url: '/assets/img/defaults/featured.png',
        thumb: '/assets/img/defaults/featured.png',
    },
    headline: {
        type: 'image',
        url: '/assets/img/defaults/featured.png',
        thumb: '/assets/img/defaults/featured.png',
    },
    streamURL: '',
    campaigns: generateDefaultCampaigns,
    supporters: [],
}

let mod = {
    defaults,
    feed: {
        load: loadFilteredFeed,
    },
    user: {
        current: getCurrentUser,
        setAvatar: setAvatar,
        updateBasicPageDetails,
        publishPage,
        saveBasics,
        updateNotifications,
        updatePageConfigurables,
        loggedIn: isUserLoggedIn,
        logout: userLogout,
    },
    public: {
        getCreator,
        getCampaign : publicGetCampaign,
        getCampaigns: publicGetCampaigns,
        getMetrics,
        pushMetrics,
        requestFileDownloadLink,
        requestService,
        getServiceRequest,
        getSiteConfig,
    },
    summary: {
        latest: getLatestSummary,
    },
    supporters: {
        get: getSupporter,
        recent: getRecentSupportersList,
        recent_campaign: getRecentCampaignSupportersList,
        hideActivity: hideSupportActivity,
        modifyServiceOrderStatus,
        all: (cached) => {return getRecentSupportersList(90000,0, cached)}
    },
    subscriptions: {
        get: getSubscription,
        listAll: listAllSubscriptions,
    },
    files: {
        delete: deleteFiles,
    },
    wallet: {
        summary: getWalletSummary,
        getPayment: getPayment,
        recent_operations: getRecentWalletOperations,
        recent_payments: getRecentPayments,
        withdraw: sendWithdrawalRequest,
        updateCashout: updateCashout,
        payoutDetails: payoutDetails,
    },
    campaign: {
        listAll: getAllCampaigns,
        get: getCampaign,
        update: updateCampaign,
        createService,
        createCampaign,
        deleteCampaign,
    },
    config: {
        serviceTemplates,
        uploadMediaEndpoint,
        gaTrackingCode: 'G-RDLPSFK3YD',
    },
    assets: {
        profPicURL,
        imageURL,
    },
    payments: {
        pollPayment,
        dispatchPayment,
        unlockCampaign,
        syncCurrenyRate,
    },
    security: {
        login: loginUser,
        signup: signupUser,
        resendSMSCode,
        resendEmailCode,
        resetPasswordRequest,
        processPasswordReset,
        verifyPhone,
        verifyEmail,
        verifyByPayment,
        updatePhone,
        updateEmail,
        updatePassword,
    },
    util: {
        validateEmail,
        validatePhone,
    },
    page: {
        set: (page) => {
            document.title = page.title + ' - MyHustle'
        
        },
        track: async () => {
            ReactGA.pageview(window.location.href)
        },
        event: (category, action, label) => {
            ReactGA.event({
                category,
                action,
                label,
                nonInteraction: false,
            })
        },
        sysEvent: (category, action, label) => {
            ReactGA.event({
                category,
                action,
                label,
                nonInteraction: true,
            })
        },
}
}

export default mod;