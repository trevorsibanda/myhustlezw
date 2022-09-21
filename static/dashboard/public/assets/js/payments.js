
function payEcocash() {
    // instanciate new modal
    var modal = new tingle.modal({
        footer: true,
        stickyFooter: false,
        closeMethods: ['overlay', 'button', 'escape'],
        closeLabel: "Close",
        cssClass: ['custom-class-1', 'custom-class-2'],
        onOpen: function () {
            console.log('modal open');
        },
        onClose: function () {
            console.log('modal closed');
        },
        beforeClose: function () {
            // here's goes some logic
            // e.g. save content before closing the modal
            return true; // close the modal
            return false; // nothing happens
        }
    });

    // set content

    modal.setContent(`
<div class="row">
    <div class="col-12">

<ul class="nav bg-light nav-pills rounded nav-fill mb-3" role="tablist">
	<li class="nav-item">
		<a class="nav-link active" data-toggle="pill" href="#nav-tab-ecocash">
		<i class="fa fa-credit-card"></i> Ecocash</a></li>
	<li class="nav-item">
		<a class="nav-link" data-toggle="pill" href="#nav-tab-cc">
		<i class="fab fa-paypal"></i>  Credit card</a></li>
</ul>
<div class="tab-content">
<div class="tab-pane fade show active" id="nav-tab-card">
	<p class="alert alert-success">Some text success or error</p>
	<form role="form">
    <div class="form-group">
        <label>* Ecocash number</label>
        <div class="input-group mb-3">
            <div class="input-group-prepend">
                <span class="input-group-text"><i class="fa fa-phone"></i></span>
            </div>
            <input type="tel" class="form-control" placeholder="07(7|8)" value="{{fanPhone}}"/>
        </div>
    </div>

	<div class="form-group">
		<label for="cardNumber">Send payment confirmation SMS</label>
		<div class="input-group">
            <select class="form-control" name="confirmationSMS">
                <option>Yes, send to ecocash number</option>
                <option>No, dont send confirmation</option>
            </select>
		</div>
	</div> <!-- form-group.// -->

</div> <!-- tab-pane.// -->
<div class="tab-pane fade" id="nav-tab-cc">
<p>Paypal is easiest way to pay online</p>
<p>
<button type="button" class="btn btn-primary"> <i class="fab fa-paypal"></i> Log in my Paypal </button>
</p>
<p><strong>Note:</strong> Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod
tempor incididunt ut labore et dolore magna aliqua. </p>
</div>
</div> <!-- tab-content .// -->
 
    </div>
</div>
`);


    let paid_thank_you = `
<div class="row">
    <div class="col-12">
        <h4>Payment received :) Thank you</h4>
        <div class="row justify-content-center">
            <div class="col-6">
                <img src="/assets/img/ecocash_paid.gif" class="img-responsive" />
            </div>
        </div>
        <strong>Your payment was successfully completed. Thank you for the support.</strong>
    </div>
</div>
`



    let lastStatus = 'queued'
    let cancelFunc = null
    let channel = null
    let pusher = null;
    let counter = 0;

    let paymentReceived = (data) => {
        lastStatus = 'paid'
        modal.setContent(paid_thank_you)
        channel.close()
        modal.addFooterBtn('Close window', 'tingle-btn tingle-btn--danger', function () {
            //create
            modal.close();
        });
        clearInterval(cancelFunc)
        pusher.disconnect()
        channel.close()

    }

    let paymentSent = (data) => {
        lastStatus = 'sent'
        modal.setContent(`
<div class="row">
    <div class="col-12">
        <h4>Authorizing payment</h4>
        <div class="row justify-content-center">
            <div class="col-6">
                <img src="/assets/img/loading.gif" class="img-responsive" />
            </div>
        </div>
        <strong>Keep your eyes on your phone! A authorization prompt should appear for you to pay</strong>
        <p>
        <small>Payment ID: <b>`+ data._id + `</b></small>
        </p>
    </div>
</div>
                `)
    }

    let paymentCancelled = (data) => {
        lastStatus = 'cancelled'
        modal.setFooterContent(``)
        modal.setContent(`
<div class="row">
    <div class="col-12">
        <h4>Payment failed :(</h4>
        <div class="row justify-content-center">
            <div class="col-6">
                <img src="/assets/img/payment_failed.gif" class="img-responsive" />
            </div>
        </div>
        <strong>`+ data.reason + `</strong>
    </div>
</div>
                `)
        pusher.disconnect()
        clearInterval(cancelFunc)
        channel.close()
    }


    let pollPayment = (id, ts, signature) => {
        if (lastStatus == 'paid' || lastStatus == 'cancelled') {
            clearInterval(cancelFunc)
            return
        }
        $.get('/api/v1/public/transaction/xhr_poll/' + id + '?ts=' + ts + '&signature=' + signature, (data, status) => {
            if (!data.status || data.status == lastStatus) {
                counter = counter + 1
                if (counter > 20) {
                    paymentCancelled({ reason: "Status has not beign updated for 300 seconds. Assuming payment failed." })
                }
                return
            }
            switch (data.status) {
                case 'paid': paymentReceived(data); break;
                case 'cancelled': paymentCancelled(data); break;
                case 'sent': paymentSent(data); break;
                default: console.log('Unknown transaction service response: ', data)
            }
        })
    }

    // add a button
    modal.addFooterBtn('Pay with Ecocash now', 'tingle-btn tingle-btn--primary', function () {

        let data = {
            items: 1,
            fullname: $('#modal_fullname').val(),
            email: $('#modal_email').val(),
            phone: $('#modal_ecocash_msidn').val(),
            message: $('#modal_message').val(),
            type: 'ecocash',
            gateway: 'ecocash',
            creator: "{{ .Creator.Username }}",
            Campaign: "{{ .Campaign.ID.Hex }}",
            purpose: 'support_creator'
        }
        modal.setFooterContent(`<small>Payment in progress. Do not leave this window.</small>`)
        $.post("/api/v1/public/transaction/initiate/support", JSON.stringify(data), (resp, status) => {
            pusher = new Pusher('d3c8fc4f4d1b77ff0011', {
                cluster: 'eu'
            });


            if (resp._id) {

                //also do http polling every 10seconds
                cancelFunc = setInterval(_ => pollPayment(resp._id, resp.ts, resp.signature), 10000)

                paymentSent(resp)

                channel = pusher.subscribe(resp._id);
                channel.bind('sent', paymentSent);

                channel.bind('cancelled', paymentCancelled);

                channel.bind('paid', paymentReceived);

            } else {
                modal.setContent('<h1>An error occured :(</h1><p class="padding-top-40">' + resp.error + '</p>')

            }
        }).catch(err => paymentCancelled({
            reason: "The transaction was rejected by the server\n" + err
        }))

    });

    // add another button
    modal.addFooterBtn('International payments/Mastercard', 'tingle-btn tingle-btn--danger', function () {
        //create
        modal.close();
    });

    // open modal
    modal.open();
}


function ecocashTransaction(msidn, amount) {
    let payload = {
        phone: msidn,
        amount: amount
    }
    return payload;
}
