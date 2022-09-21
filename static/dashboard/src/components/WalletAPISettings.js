import {Component} from "react"
import v1 from "../api/v1"


class WalletAPISettings extends Component {
    constructor(props){
        super(props)

        this.state = {
            api: {
                enable: false,
                url: "",
                method: "",
            }
        }
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
                <h4 class="box-title mb-15">Wallet API Settings</h4>
                <div class="pad">
                    <div class="row">
                        <div class="col-lg-7 col-md-6 col-12">
                            <div class="row">
                                <div class="col-5">
                                    <div class="form-group">
                                        <label>Target URL</label>
                                        <input type="url" class="form-control" placeholder="https://" />
                                    </div>
                                </div>
                                <div class="col-7 pull-right">
                                    <div class="form-group">
                                        <label>Method</label>
                                                <select class="form-control" value={active.bankname} onChange={(evt) => this.updateUIField('bankname', evt)}  >
                                                    <option value="post_json" >POST JSON</option>
                                                    <option value="post_form">POST Form</option>
                                                    <option value="get" >GET</option>
                                        </select>
                                    </div>
                                </div>
                            </div>

                            
                            <button class="btn btn-success btn-rounded" onClick={this.doUpdateDetails} >Save details</button>
                        </div>
                        <div class="col-lg-5 col-md-6 col-12">
                            <h3 class="box-title mt-10">API Config</h3>
                            <p>We will send a request on all new payments to the configured URL. </p>
                        </div>
                    </div>
                </div>
            </>
        )
    }
}


export default WalletAPISettings;
