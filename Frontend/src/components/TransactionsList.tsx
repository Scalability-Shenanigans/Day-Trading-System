import React, { useState } from "react";
import styled from "styled-components";
import { commitBuy, commitSell } from "../requests/requests";

const TransactionsListContainer = styled.div`
  flex: 1;
  box-sizing: border-box;
  overflow-y: auto;
  width: 100%;
  background: #2c3e50;
  border-radius: 5px;
  padding: 10px;
`;

const TransactionsTable = styled.table`
  width: 100%;
  border-collapse: collapse;
  box-sizing: border-box;
  color: #ecf0f1;
`;

const TransactionsTableHeader = styled.thead`
  font-weight: bold;
  box-sizing: border-box;
`;

const TransactionsTableRow = styled.tr`
  background: #34495e;
  border-radius: 3px;
  &:nth-child(even) {
    background: #2c3e50;
  }
  box-sizing: border-box;
`;

const TransactionsTableCell = styled.td`
  padding: 8px;
  text-align: center;
  box-sizing: border-box;
  border-bottom: 1px solid #7f8c8d;
`;

export interface TransactionsListItemProps {
  type: string;
  date: string;
  asset: string;
  amount: number;
  isCommitted: boolean;
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
  isCommitted,
}) => {
  return (
    <TransactionsTableRow>
      <TransactionsTableCell>{type}</TransactionsTableCell>
      <TransactionsTableCell>{date}</TransactionsTableCell>
      <TransactionsTableCell>{asset}</TransactionsTableCell>
      <TransactionsTableCell>{amount}</TransactionsTableCell>
      {isCommitted === true ? (
        <TransactionsTableCell style={{ color: "lightgreen" }}>
          Committed
        </TransactionsTableCell>
      ) : (
        <TransactionsTableCell style={{ color: "orange" }}>
          Not Committed
        </TransactionsTableCell>
      )}
    </TransactionsTableRow>
  );
};

const TransactionsList: React.FC<TransactionsListProps> = ({
  transactions,
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
              isCommitted={transaction.isCommitted}
            />
          ))}
        </tbody>
      </TransactionsTable>
    </TransactionsListContainer>
  );
};

export default TransactionsList;
