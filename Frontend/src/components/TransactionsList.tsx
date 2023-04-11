import React, { useState } from "react";
import styled from "styled-components";
import { commitBuy, commitSell } from "../requests/requests";

const TransactionsListContainer = styled.div`
  box-sizing: border-box;
  min-height: 200px;
  width: 100%;
  background: lightgray;
  border-radius: 5px;
  padding: 10px;
`;

const TransactionsTable = styled.table`
  width: 100%;
  border-collapse: collapse;
`;

const TransactionsTableHeader = styled.thead`
  font-weight: bold;
`;

const TransactionsTableRow = styled.tr`
  background: white;
  border-radius: 3px;
  &:nth-child(even) {
    background: #f2f2f2;
  }
`;

const TransactionsTableCell = styled.td`
  padding: 8px;
  text-align: center;
`;

const StatusButton = styled.button`
  margin: 0;
`;

export interface TransactionsListItemProps {
  type: string;
  date: string;
  asset: string;
  amount: number;
  user: string;
  setCommit: React.Dispatch<React.SetStateAction<boolean>>;
}

interface TransactionsListProps {
  transactions: TransactionsListItemProps[];
  user: string;
  setCommit: React.Dispatch<React.SetStateAction<boolean>>;
}

const TransactionsListItem: React.FC<TransactionsListItemProps> = ({
  type,
  date,
  asset,
  amount,
  user,
  setCommit,
}) => {
  const [status, setStatus] = useState("Not committed");

  return (
    <TransactionsTableRow>
      <TransactionsTableCell>{type}</TransactionsTableCell>
      <TransactionsTableCell>{date}</TransactionsTableCell>
      <TransactionsTableCell>{asset}</TransactionsTableCell>
      <TransactionsTableCell>{amount}</TransactionsTableCell>
      {status === "Committed" ? (
        <TransactionsTableCell>Committed</TransactionsTableCell>
      ) : (
        <TransactionsTableCell>
          <StatusButton
            onClick={async () => {
              // send request here
              if (type === "Buy") {
                const commitBuyResponse = await commitBuy({
                  user: user,
                });

                if (
                  commitBuyResponse?.status === 200 &&
                  commitBuyResponse.data["stock"] === asset
                ) {
                  setStatus("Committed");
                  setCommit(true);
                }
              } else if (type === "Sell") {
                const commitSellResponse = await commitSell({
                  user: user,
                });

                if (
                  commitSellResponse?.status === 200 &&
                  commitSellResponse.data["stock"] === asset
                ) {
                  setStatus("Committed");
                  setCommit(true);
                }
              }
            }}
          >
            commit
          </StatusButton>
        </TransactionsTableCell>
      )}
    </TransactionsTableRow>
  );
};

const TransactionsList: React.FC<TransactionsListProps> = ({
  transactions,
  user,
  setCommit,
}) => {
  return (
    <TransactionsListContainer>
      <TransactionsTable>
        <TransactionsTableHeader>
          <tr>
            <TransactionsTableCell>Type</TransactionsTableCell>
            <TransactionsTableCell>Date</TransactionsTableCell>
            <TransactionsTableCell>Asset</TransactionsTableCell>
            <TransactionsTableCell>Amount</TransactionsTableCell>
            <TransactionsTableCell>Status</TransactionsTableCell>
          </tr>
        </TransactionsTableHeader>
        <tbody>
          {transactions.map((transaction, index) => (
            <TransactionsListItem
              key={index}
              type={transaction.type}
              date={transaction.date}
              asset={transaction.asset}
              amount={transaction.amount}
              user={user}
              setCommit={setCommit}
            />
          ))}
        </tbody>
      </TransactionsTable>
    </TransactionsListContainer>
  );
};

export default TransactionsList;
