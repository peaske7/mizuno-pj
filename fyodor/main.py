import torch
from PIL import Image
from transformers import AutoModel, AutoTokenizer, pipeline, AutoModelForCausalLM

def cpm():
    model = AutoModel.from_pretrained(
        "openbmb/MiniCPM-Llama3-V-2_5",
        trust_remote_code=True,
        torch_dtype=torch.float16,
    )

    tokenizer = AutoTokenizer.from_pretrained(
        "openbmb/MiniCPM-Llama3-V-2_5", trust_remote_code=True
    )
    model.eval()

    image = Image.open(".img/image-00000.jpg").convert("RGB")
    question = "What is in the image?"
    msgs = [{"role": "user", "content": question}]

    res = model.chat(
        image=image,
        msgs=msgs,
        tokenizer=tokenizer,
        sampling=True,  # if sampling=False, beam_search will be used by default
        temperature=0.7,
        # system_prompt='' # pass system_prompt if needed
    )
    print(res)

    ## if you want to use streaming, please make sure sampling=True and stream=True
    ## the model.chat will return a generator
    res = model.chat(
        image=image,
        msgs=msgs,
        tokenizer=tokenizer,
        sampling=True,
        temperature=0.7,
        stream=True,
    )

    generated_text = ""
    for new_text in res:
        generated_text += new_text
        print(new_text, flush=True, end="")


def layout():
    nlp = pipeline(
        "document-question-answering",
        model="impira/layoutlm-document-qa",
    )

    invoice_number = nlp(
        ".img/image-00000.jpg",
        "What is the 伝票NO?",
    )

    received_at = nlp(
        ".img/image-00000.jpg",
        "What is the 入庫日?",
    )

    who = nlp(
        ".img/image-00000.jpg",
        "Who is the sender (the name before 御中?)",
    )

    items = nlp(
        ".img/image-00000.jpg",
        # "Output the list of 商品コード, 商品名, 荷姿, 単量, 入庫数, 備考 in csv format",
        "What is the 商品コード",
    )

    print(invoice_number)
    print(received_at)
    print(who)
    print(items)

def monkey():
    checkpoint = "echo840/Monkey-Chat"
    model = AutoModelForCausalLM.from_pretrained(checkpoint, trust_remote_code=True).eval()
    tokenizer = AutoTokenizer.from_pretrained(checkpoint, trust_remote_code=True)
    tokenizer.padding_side = 'left'
    tokenizer.pad_token_id = tokenizer.eod_id
    img_path = ".img/image-00000.jpg"
    question = "Output the list of 商品コード, 商品名, 荷姿, 単量, 入庫数, 備考 in csv format"

    #Monkey-Chat has the same prompt format for both vqa and detailed caption.
    query = f'<img>{img_path}</img> {question} Answer: '

    input_ids = tokenizer(query, return_tensors='pt', padding='longest')
    attention_mask = input_ids.attention_mask
    input_ids = input_ids.input_ids

    pred = model.generate(
                input_ids=input_ids,
                attention_mask=attention_mask,
                # do_sample=False,
                num_beams=1,
                max_new_tokens=512,
                min_new_tokens=1,
                length_penalty=1,
                num_return_sequences=1,
                output_hidden_states=True,
                use_cache=True,
                pad_token_id=tokenizer.eod_id,
                eos_token_id=tokenizer.eod_id,
                )
    response = tokenizer.decode(pred[0][input_ids.size(1):].cpu(), skip_special_tokens=True).strip()
    print(response)

if __name__ == "__main__":
    # cpm()
    # layout()
    monkey()
