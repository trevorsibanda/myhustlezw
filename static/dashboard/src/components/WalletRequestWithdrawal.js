import {Component} from "react"
import v1 from "../api/v1"
import money from "./payments/Amount"

class WalletRequestWithdrawal extends Component {
    constructor(props){
        super(props)

        this.state = {
            processing: this.props.summary.pending_withdrawal > 0,
        }

        this.doWithdraw = this.doWithdraw.bind(this)
    }


    doWithdraw(){
        let currency = this.props.currency
        v1.wallet.withdraw(currency).then(resp => {
            if(resp.status === 'ok'){
                this.setState({processing: true})
                alert("We have received your withdrawal request. You will receive an email and SMS when your request has been processed.")
                return this.props.onWithdraw ? this.props.onWithdraw(currency, resp.summary): 0
            }else{
                alert("An error occured whilst processing your request!")
            }
        }).catch(err => {
           
        })
    }

    render() {
        return ( this.state.withdrawn || this.props.summary.pending_withdrawal > 0.00 ? <>
            <h4 class="box-title mb-15">Your withdrawal request has been received</h4>
            <p>To see which bank the payout details will be used, please check your email.</p>
        </> :
            <>
                <h4 class="box-title mb-15">{this.props.currency} Withdrawal Request</h4>
                <div class="pad">
                    <div class="row">
                        <div class="col-lg-7 col-md-6 col-12">
                            <div class="form-group">
                                <label for="exampleInputEmail1">Amount to withdraw</label>
                                <div class="input-group">
                                    <div class="input-group-addon">{this.props.currency} $</div>
                                    <input type="number" class="form-control" placeholder="0.00" max={this.props.summary.available} value={this.props.summary.available} readonly min={this.props.min_withdraw} />
                                </div>
                                <small>Minimum withdrawal amount is {money.format(this.props.currency, this.props.min_withdraw)} <a href="javascript:;" onClick={this.requestWithdraw} style={{ "color": "red" }}>Withdraw maximum amount</a></small>
                            </div>
                            <div class="alert alert-primary">
                                <p>MyHustle charges a flat fee of <b>30%</b>. This means you will receive 70% of the amount shown above in your bank account.</p>
                                <p>You have {money.format(this.props.currency, this.props.summary.available)} and you will receive a payout
                                    of {money.format(this.props.currency, this.props.summary.available * 0.7)}
                                    and we will retain {money.format(this.props.currency, this.props.summary.available * 0.3)}</p>
                            </div>
                            <button class="btn btn-success btn-rounded" disabled={this.props.summary.available < this.props.min_withdrawal} onClick={this.doWithdraw} ><i class="fa fa-document"></i> Request withdrawal</button>
                        </div>
                        <div class="col-lg-5 col-md-6 col-12">
                            <h3 class="box-title mt-10">Withdrawal process</h3>

                            <p>Withdrawals are processed twice a week, on Tuesday and Thursday with a cutoff time of 12PM GMT+2.
                                Any withdrawal requests after then are not guaranteed to be effected</p>
                            <p>Reminder: You have two wallets, you are currently viewing the {this.props.currency} wallet. You can switch wallets at the top of your page.</p>
                        </div>
                    </div>
                </div>
                <hr />
            </>
        )
    }
}

export default WalletRequestWithdrawal;
