const faqs = [
  {
    question: "Is there a free version?",
    answer: "Yes!",
  },
  {
    question: "What are your support hours?",
    answer:
      "We have support staff answering emails 24/7, though response times may be a bit slower on weekends.",
  },
  {
    question: "How do I contact support?",
    answer: (
      <>
        Email us -{" "}
        <a href="mailto:support@lenslocked.com">support@lenslocked.com</a>{" "}
      </>
    ),
  },
];

export default function Faq() {
  return (
    <>
      <div className="px-6">
        <h1 className="py-4 text-4xl semibold tracking-tight">FAQ Page</h1>
        <ul className="grid grid-cols-2 gap-16">
          {faqs.map((faq) => (
            <QAPair question={faq.question} answer={faq.answer} />
          ))}
        </ul>
      </div>
    </>
  );
}

interface QAPairProps {
  question: string;
  answer: JSX.Element | string;
}

function QAPair({ question, answer }: QAPairProps) {
  return (
    <li className="my-4 border-t border-indigo-400 py-1 px-2">
      <span className="block text-lg text-gray-800 semibold">{question}</span>
      <span className="block text-sm text-gray-500">{answer}</span>
    </li>
  );
}
