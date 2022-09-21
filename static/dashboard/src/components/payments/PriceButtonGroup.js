

function PriceButtonGroup(props) {
    return (
        <div class="btn-group btn-group-justified" role="group" aria-label="...">
            {props.prices.map((price) => {
                return (
                    <div class="btn-group" role="group">
                        <button type="button" class={"btn "+ (props.price === price ? "btn-info" : "btn-default")} onClick={() => props.onChange(price)}>$ {price}</button>
                    </div>
                )
            })}
        </div>
    )
}

export default PriceButtonGroup;