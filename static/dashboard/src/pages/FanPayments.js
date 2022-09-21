import { Component } from "react"
import { Link } from "react-router-dom";
import ReactTimeAgo from "react-time-ago/commonjs/ReactTimeAgo";
import v1 from "../api/v1";



class FanPayments extends Component {
    constructor(props) {
        super(props)

        this.state = {
            payments: []
        }

        v1.wallet.recent_payments(10, false).then(operations => {
            this.setState({ payments: operations })
        }).catch(err => {
            alert('Failed to fetch recent payments\n\nError:' + err)
        })

    }

    supportType(payment) {
        let s = "Unknown"
        switch(payment.action) {
            case "verify_account": s = "Account Verification Fee"; break;
            case "add_supporter_message": s = "Bought someone "+ payment.items + "(s) " + payment.item_name + "(s)"; break;
            case "grant_service": s = "Paid for a service"; break;
            case "grant_subscribe": s = <>Subscribed to <Link to={"/@"+payment.item_name }>@{payment.item_name}</Link></>; break;
            case "grant_campaign": s = "Paid to view content"; break;
            default:
                s = "Unknown payment purpose";
                break;
        }
        return s;
    }


    render() {
        let paymentList = this.state.payments.map(payment => {
            return <div class="col-md-6 col-12">
                <div class="box">
                    <div class="box-header with-border">
                        <h4 class="box-title box-title-bold">{this.supportType(payment)}</h4>
                        <h5 class="box-title ">{payment.currency} {payment.price} </h5>
                    </div>
                    <div class="box-body">
                        <p>Paid <b>{payment.currency} {payment.price}</b> using <b>{payment.gateway}</b> <ReactTimeAgo  date={payment.created} />.</p>
                        <p>PaymentID: <strong>{payment._id}</strong> at {payment.created}</p>
                    </div>
                </div>
            </div>
        })

        return (
            <>
                <h4>My payments</h4>
                <div class="row justify-content-center">
                      {paymentList}
                </div>
                
            </>
        )
    }
}


export default FanPayments;
