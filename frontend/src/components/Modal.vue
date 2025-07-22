<template>
    <div class="modal">
        <div class="modal-content">
            <div class="modal-header">
                <h2>
                    <slot name="title"></slot>
                </h2>
                <i class="uil uil-times close" @click="$emit('closeModal')"></i>
            </div>
            <slot name="body"></slot>
        </div>
    </div>
</template>


<style>
.modal {
    position: fixed;
    height: 100vh;
    width: 100vw;
    top: 0;
    left: 0;
    background-color: rgba(0, 0, 0, 0.4);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 2;

}

.modal-content {
    display: flex;
    flex-direction: column;
    gap: 25px;

    padding: 30px 50px;
    width: 400px;



    position: relative;
    background-color: var(--page-bg);
    border-radius: var(--container-border-radius);


}

.modal-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.modal-header i {
    font-size: 1.25em;
    cursor: pointer;

}

.modal-content form {
    display: flex;
    flex-direction: column;
    gap: 15px;
}

.modal-content .btn.form-submit {
    align-self: flex-end;
}

/* close icon */
.uil-times.close {
    cursor: pointer;
    transition: transform 0.25s ease;
}

.uil-times.close:hover {
    transform: scale(1.15);
}
</style>