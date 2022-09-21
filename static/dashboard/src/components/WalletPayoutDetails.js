import {Component} from "react"
import v1 from "../api/v1"


class WalletUpdatePayoutDetails extends Component {
    constructor(props){
        super(props)

        this.state = {
            currency: this.props.currency,
            ZWL: {
                method: "bank",
                bankname: "",
                bankbranch: "",
                bankaccountname: "",
                bankaccountnumber: "",
                phonenumber: "",
            },
            USD: {
                method: "bank",
                bankname: "",
                bankbranch: "",
                bankaccountname: "",
                bankaccountnumber: "",
                phonenumber: "",
            }
        }

        v1.wallet.payoutDetails(false).then(resp => {
            this.setState({
                ZWL: resp.ZWL ? resp.ZWL : this.state.ZWL,
                USD: resp.USD ? resp.USD : this.state.USD,
            })
        }).catch(err => {
            console.log('failed to fetch payout details', err)
        })

        this.supportedBanks = ['AgriBank', 'BancABC', 'CABS', 'CBZ', 
        'EcoBank', 'FBC Bank', 'First Capital', 'NedBank', 'NMB', 
        'NBS', 'POSB', 'Stanbic', 'Standard Chartered', 'ZB Bank',
        'Empower Bank', 'FBC Bld Society', 'African Century', 'Success Bank', 'GetBucks', 'MetBank', 'LION', 'ZWMB']

        this.updateUIField = (field, evt) => {
            let active = this.state[this.props.currency]
            active[field] = evt.target.value
            if(this.props.currency === "ZWL"){
                this.setState({ZWL: active})
            }else{
                this.setState({USD: active})
            }
        }


        this.doUpdateDetails = this.doUpdateDetails.bind(this)
    }

    doUpdateDetails(){
        v1.wallet.updateCashout(this.props.currency, this.state[this.props.currency]).then(resp => {
            if(resp.status === 'ok'){
                alert("Successfully updated your payout details.")
            }else{
                alert("Failed to update payout details:\n\nReason:"+resp.error)
            }
        }).catch(err => {
            alert("An error occured. Error:"+ err)
        })
    }

    render() {
        let active = this.state[this.props.currency]
        return (
            <>
                <h4 class="box-title mb-15">{this.props.currency} Bank/Payout Details</h4>
                <div class="pad">
                    <div class="row">
                        <div class="col-lg-7 col-md-6 col-12">
                            <div class="row">
                                <div class="col-5">
                                    <div class="form-group">
                                        <label>Withdrawal method</label>
                                        <select class="form-control" value={active.method} onChange={(evt) => this.updateUIField('method', evt) }  >
                                            <option value="bank">Bank transfer</option>
                                            <option value="mobile_money">Mobile money</option>
                                        </select>
                                    </div>
                                </div>
                                <div class="col-7 pull-right">
                                {
                                active.method === "bank" ?
                                    <div class="form-group">
                                        <label>Bank name</label>
                                                <select class="form-control" value={active.bankname} onChange={(evt) => this.updateUIField('bankname', evt)}  >
                                            {this.supportedBanks.map( bank => {
                                                return <option value={bank}>{bank}</option>
                                            })}
                                        </select>
                                    </div> :
                                    <div class="form-group">
                                        <label for="exampleInputEmail1">Mobile Number</label>
                                        <div class="input-group">
                                            <div class="input-group-addon">+263</div>
                                                    <input type="text" class="form-control" value={active.phonenumber} onChange={(evt) => this.updateUIField('phonenumber', evt) } placeholder="Phone number" />
                                        </div>
                                    </div>
                                }
                                </div>
                            </div>
                            {active.method === "bank" ? 
                            <>
                            <div class="form-group">
                                <label for="exampleInputEmail1">Branch name</label>
                                <div class="input-group">
                                    <div class="input-group-addon"><i class="fa fa-building"></i></div>
                                            <input type="text" class="form-control" placeholder="Branch" onChange={(evt) => this.updateUIField('bankbranch', evt)} value={active.bankbranch} />
                                </div>
                            </div>
                            <div class="form-group">
                                <label for="exampleInputEmail1">Account name</label>
                                <div class="input-group">
                                    <div class="input-group-addon"><i class="fa fa-user"></i></div>
                                            <input type="text" class="form-control" placeholder="Account name" value={active.bankaccountname} onChange={(evt) => this.updateUIField('bankaccountname', evt)} />
                                </div>
                            </div>


                            <div class="row">
                                <div class="col-12">
                                    <div class="form-group">
                                        <label>Account number</label>
                                                <input type="text" class="form-control" value={active.bankaccountnumber} onChange={(evt) => this.updateUIField('bankaccountnumber', evt)} placeholder="Account number" />
                                    </div>
                                </div>
                            </div>
                            </>
                            : <><p>Payouts only available to Ecocash numbers</p></>}

                            
                            <button class="btn btn-success btn-rounded" onClick={this.doUpdateDetails} >Save details</button>
                        </div>
                        <div class="col-lg-5 col-md-6 col-12">
                            <h3 class="box-title mt-10">Need help?</h3>
                            <p>Please reach out using anyone of the contacts listed on the contact page. </p>
                        </div>
                    </div>
                </div>
            </>
        )
    }
}


export default WalletUpdatePayoutDetails;
