import { Link } from "react-router-dom";
import money from "./payments/Amount";

function PageEarnings(props) {
    return (
                <><div class="row">
                    <div class="col-lg-3 col-md-6">
                <div class={"box box-inverse" + (props.earnings.zwl > 0 ? " bg-success" : "bg-gradient-botani")}>
                            <div class="box-body">
                                <h3>{money.formatZWL(props.earnings.zwl)} </h3>
                                <span class="pull-right">earned through this page</span>
                            </div>
                        </div>
                    </div>
                    <div class="col-lg-3 col-md-6">
                        <div class={"box box-inverse" + (props.earnings.usd > 0 ? " bg-success" : " bg-gradient-botani")}>
                            <div class="box-body">
                                <h3>{money.formatUSD(props.earnings.usd)}</h3>
                                <span class="pull-right">earned through this page</span>
                            </div>
                        </div>
                    </div>
        </div>
            <div class="row mt-20" >
                <Link to="/creator/wallet" class="btn btn-info btn-block"><i class="fa fa-credit-card"></i> View Wallet</Link>
            </div>
            </>
    )
}

export default PageEarnings;