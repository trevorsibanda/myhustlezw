import React, { Component } from "react";
import v1 from "../api/v1";
import money from "../components/payments/Amount"
import { Link } from "react-router-dom"


import ReactTimeAgo from 'react-time-ago'
import Preloader from "../components/PreLoader";

class ViewPayment extends Component {

    constructor(props) {
        super(props)

        this.state = {
            loaded: false,
            payment: {
                _id: 'Loading',
                amount: 0,
                currency: 'Loading',
                created_at: 'Loading',
                expires: 'Loading',
            }
        }

        let { id } = this.props.match.params


        v1.payments.get(id).then(payment => {
            this.setState({ payment, loaded: true })
        }).catch(_ => {
            this.setState({loaded: false})
        })
    }

    render() {
        return this.state.loaded ?  (
            <>
                <div class="box">
                    <div class="box-header">
                        <div class="row">
                            <div class="col-md-8">
                                <h4 class="box-title">Payment: {this.state.payment._id}</h4>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="media-body">
                    <table class="table table-responsive">
                        <tbody>
                            <tr>
                                <td>Paid</td>
                                <td>{money.format(this.state.payment.currency, this.state.payment.amount)}</td>
                            </tr>
                            <tr>
                                <td>Subscribed</td>
                                <td><ReactTimeAgo date={this.state.payment.created_at} /></td>
                            </tr>
                            <tr>
                                <td>Expires</td>
                                <td><ReactTimeAgo date={this.state.payment.expires} /></td>
                            </tr>
                        </tbody>
                    </table>

                </div>

            </>
        ) : <Preloader />
    }

}

export default ViewPayment;