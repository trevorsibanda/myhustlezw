import {Component} from "react"
import v1 from "../api/v1";
import money from "./payments/Amount"


class RecentWalletOperations extends Component {
    constructor(props) {
        super(props)

        this.state = {
            events: [{
                operation: 'loading',
                created_at: '',
                currency: 'ZWL',
                amount: 0.00,
                gateway: 'Paynow',
            }]
        }

        v1.wallet.recent_operations(10, false).then(operations => {
            this.setState({ events: operations })
        }).catch( _ => {
            v1.wallet.recent_operations(10, true).then(operations => {
                this.setState({ events: operations })
            })
        })

    }


    render() {
        let eventList = this.state.events.map(event => {
            return <tr>
                <th scope="row">{event.operation}</th>
                <td>{event.currency}</td>
                <td>{money.format(event.currency, event.amount)}</td>
                <td>{event.gateway}</td>
                <td>{event.last_updated}</td>
            </tr>
        })

        return (
            <>
                <h4>Pending operations</h4>
                <div class="table-responsive">
                    <table class="table mb-0">
                        <thead class="thead-dark">
                            <tr>
                                <th scope="col">Event</th>
                                <th scope="col">Currency</th>
                                <th scope="col">Amount</th>
                                <th scope="col">Gateway</th>
                                <th scope="col">When</th>
                            </tr>
                        </thead>
                        <tbody>
                            {eventList}

                        </tbody>
                    </table>
                </div>
            </>
        )
    }
}


export default RecentWalletOperations;
